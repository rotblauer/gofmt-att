package fmtatt

import (
	"github.com/rotblauer/gofmt-att/remote"
	"path/filepath"
	"time"
	"github.com/kr/pretty"
)

/*
workerLoop handles
- clone
- fmt
- iff dirty
	- branch and commit
	- spawn contributorFn function
  else
	- teardown
*/




func (f *FmtAtt) processRepo(r *remote.RepoT, outcome *remote.Outcome) (moarRepos []*remote.RepoT, err error) {
	if f.Config.GitConfig.BasePath == "" {
		panic("empty git config base path")
	}
	r.Target = filepath.Join(f.Config.GitConfig.BasePath, r.Owner.Name, r.Name)

	f.Logger.I("working", r)


	// because we don't store this data, weirdly enough
	// FIXME?
	if r.CloneUrl == "" {
		r.GuessURLs(f.Config.RepoProvider)
	}

	// clone
	f.Logger.If("cloning (%d/%d) %s", f.workingPool.len(), cloningWorkLimit, r)
	var start = time.Now()
	err = f.Giter.Clone(r)
	if err != nil {
		return
	}
	f.Logger.If("cloning finished %s (took %v)", r, time.Since(start).Round(time.Second))

	// fmt
	f.Logger.If("fmting repo %s", r)
	err = f.fmter(r)
	if err != nil {
		return
	}
	f.Logger.If("fmting finished %s (took %v)", r, time.Since(start).Round(time.Second))

	// get status
	dirty, status, err := f.Giter.Status(r.Target)
	if err != nil {
		return
	}
	if !dirty {
		outcome.Status = remote.Clean
		return
	}

	outcome.Status = remote.Dirty
	outcome.FormattedOutcome = &remote.FormattedOutcome{GitStatus: status}
	f.Logger.I(r.String(), outcome.String())

	// get all repos from stripe spec
	// this doesn't depend on the outcome of the adding spec
	// kind of seems upside-down, get over it
	moarRepos = f.stripeStatus(status)

	added, err := f.add(r, outcome, status)
	if err != nil {
		return
	}
	// if our list-ifying excluded all proposed changes, return
	if added == 0 {
		outcome.Status = remote.Clean
		return
	}

	// commit and branch
	hash, status, err := f.Giter.CommitWithBranch(r.Target, f.Config.GitConfig.GitCommitConfig)
	if err != nil {
		return
	}
	outcome.FormattedOutcome.GitStatus = status
	outcome.FormattedOutcome.CommitHash= hash
	outcome.Status = remote.Committed
	return
}

var cloneQ = make(chan int, cloningWorkLimit) // semaphore
func (f *FmtAtt) workerLoop() {
	f.Logger.I("starting worker")

	for r := range f.workerChan {
		cloneQ <- 1
		for f.pause {} // let any running workers finish up. kind of messy but whatever
		o := f.mustGetRepoOutcome(r)
		if o == nil {
			f.Logger.E("no persisted outcome", r.String())
			<-cloneQ
			continue
		}
		go func(r *remote.RepoT, outcome *remote.Outcome) {
			f.workingPool.push(r)
			stripedRepos, err := f.processRepo(r, outcome)
			f.mustPutRepoOutcome(r, outcome)
			f.Logger.I("finished processing")
			f.Logger.If("%s %s", r.String(), outcome.String())
			if err != nil {
				outcome.SetErr(err)
				f.Logger.Ef("ERR_PROCESSING=%v %s %s", err, r.String(), outcome.String())
				f.Logger.E(pretty.Sprint(r))
				f.Logger.E(pretty.Sprint(outcome))
			}
			for _, rr := range stripedRepos {
				f.Logger.If("++ %s %s", rr.String())
				f.striperChan <- rr
			}
			if outcome.Status == remote.Committed {
				f.Logger.If("sending %s to contributor channel", r.String())
				f.contributorChan <- struct{
					r *remote.RepoT
					o *remote.Outcome
				}{
					r: r,
					o: outcome,
				}
			} else {
				f.teardown(r)
			}
			<-cloneQ
			f.workingPool.splice(r)
		}(r, o)
	}
	f.Logger.W("quitting worker loop")
}