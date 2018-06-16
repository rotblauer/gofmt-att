package fmtatt

import (
	"github.com/rotblauer/gofmt-att/persist"
	"github.com/rotblauer/gofmt-att/remote"
	"time"
	"github.com/kr/pretty"
	"errors"
)

func (f *FmtAtt) filterRawRepo(r *remote.RepoT, ff interface{}) (ok bool, outcome *remote.Outcome, filterErr *remote.ErrFilteredT) {

	if f.repoPool.GetOutcome(r) != nil {
		return
		// continue // dedupe, already in pool
	}

	filterErr = f.Repoer.FilterRepo(ff, f.Config.ReposSpec.RepoListSpec)
	if filterErr != nil {
		outcome = &remote.Outcome{
			Status:           remote.Invalid,
			FilteredE:        filterErr,
		}
		return
	}
	// try to get from persister
	outcome = f.mustGetRepoOutcome(r)
	// if was nil, initialize and return ok
	if outcome == nil {
		outcome = &remote.Outcome{
			Status:           remote.Valid,
		}
		ok = true
		return
	}

	outcomeExpired := outcome.Timestamp.Before(time.Now().Add(-f.Config.ReposSpec.FmtExpiration))

	outcomeUnfinished := outcome.Status > remote.Clean // either cleared for fmting, with no progress otherwise
	outcomeUnfinished = outcomeUnfinished && outcome.Status < remote.PullRequested
	outcomeUnfinished = outcomeUnfinished && outcome.Error == "" // only retry repos that closed without an error. TODO: use a shutting down error to mark this

	// hack, fixme
	// if outcome was unfinished, we need to clear it
	// because we won't be using the same timestamped branchname
	if outcomeUnfinished {
		outcome = &remote.Outcome{
			Status:           remote.Valid,
		}
	}

	ok = outcome.Status == remote.Valid || outcomeExpired
	return
}

func (f *FmtAtt) fetch(state persist.PersistentState) {
	f.Logger.If("fetching %s", state.Current)

	fetched, err := f.Repoer.FetchLeaf(state.Current, f.Config.ReposSpec.RepoListSpec)
	if err != nil {
		f.Logger.E("ERR FETCHING", err.Error())
		return
	}
	f.Logger.If("fetched %d resources", len(fetched))

	var okRepos []*remote.RepoT
	var okOwners []*remote.Owner
	var filteredRepos = make(map[*remote.RepoT]*remote.ErrFilteredT)
	var filteredOwners = make(map[*remote.Owner]*remote.ErrFilteredT)

	for _, ff := range fetched {
		r, isRepo := f.Repoer.ToRepo(ff)
		if isRepo {
			ok, outcome, filterErr := f.filterRawRepo(r, ff)
			f.mustPutRepoOutcome(r, outcome)
			if filterErr != nil {
				filteredRepos[r] = filterErr
				continue
			}
			f.Persister.PutOwner(r.Owner)
			if !ok {
				continue
			}
			// set repo AND owner to pool
			okRepos = append(okRepos, r)
			f.repoPool.Set(r, outcome)
			if !f.ownerPool.Has(r.Owner) {
				okOwners = append(okOwners, r.Owner)
				f.ownerPool.Push(r.Owner)
			}
			continue
		}
		o, isOwner := f.Repoer.ToOwner(ff)
		if !isOwner {
			err = errors.New("failed to assert type " + pretty.Sprint(ff))
			break
		}
		filterErr := f.Repoer.FilterOwner(ff, f.Config.ReposSpec.OwnerListSpec)
		if filterErr != nil {
			filteredOwners[o] = filterErr
			continue
		}
		f.Persister.PutOwner(o)
		if !f.ownerPool.Has(o) {
			okOwners = append(okOwners, o)
			f.ownerPool.Push(o)
		}
	}
	// loggers
	if err != nil {
		// failed to cast
		f.Logger.Ef("ERR FETCHER TYPE ASSERTER = %v", err)
	}

	for o, e := range filteredOwners {
		f.Logger.Df("~ %s | %s", o.String(), e.String())
	}
	for r, e := range filteredRepos {
		f.Logger.Df("~ %s | %s", r.String(), e.String())
	}
	for _, o := range okOwners {
		f.Logger.If("+ %s", o.String())
	}
	for _, r := range okRepos {
		f.workerChan <- r // fire off to worker chan
		f.Logger.If("+ %s", r.String())
	}
}