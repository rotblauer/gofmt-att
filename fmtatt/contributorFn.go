package fmtatt

import (
	"github.com/rotblauer/gofmt-att/remote"
)

/*
contributorFn handles
- fork
- push
- pr
- teardown

It is synchronous, but spawned in a goroutine.
*/
func (f *FmtAtt) contributorFn(r *remote.RepoT, outcome *remote.Outcome) {
	if f.dryRun[0] {
		f.Logger.E("dryrun[fork], skipping fork and subsequent processes")
		return
	}

	f.Logger.I("contributing to", r.String())

	var err error

	defer func() {
		if err != nil {
			outcome.Error = err.Error()
			f.Logger.Ef("ERR: %s %s", r, outcome)
		} else {
			f.Logger.If("%s %s", r, outcome)
		}
		f.mustPutRepoOutcome(r, outcome)
		f.teardown(r)
	}()

	// sanity safety check
	if outcome != nil && outcome.Status == remote.PullRequested {
		f.Logger.W("... oops, already solved this repo. i am an idiot.")
	}

	// fork
	f.Logger.I("forking", r.String())
	of, err := f.Repoer.ForkRepo(f.Config.ForkConfig, r)
	if err != nil {
		return
	}
	// ... finished. this can take a hot sec.
	outcome.Status = remote.Forked
	outcome.ForkedOutcome = &remote.ForkedOutcome{
		GitUrl:   of[1].GitUrl,
		HTMLUrl:  of[1].HTMLUrl,
		CloneUrl: of[1].CloneUrl,
	}
	f.mustPutRepoOutcome(r, outcome) // just in case something happens
	f.Logger.I("forked", r.String(), outcome.String())

	// push
	if f.dryRun[1] {
		f.Logger.E("dryrun[push], skipping push and subsequent processes")
		return
	}

	f.Logger.I("pushing", r.String())
	err = f.Giter.PushAll(r.Target, of[1].CloneUrl, f.Config.GitConfig.BranchName) // note that this adds remote, too
	if err != nil {
		return
	}

	// pr
	if f.dryRun[2] {
		f.Logger.E("dryrun[pr], skipping pr and subsequent processes")
		return
	}
	pr := &remote.PullRequestT{
		RepoT: r,
		PullRequestConfig: &remote.PullRequestConfig{
			Title:    f.Config.PullRequestConfig.Title,
			Head:     of[1].Owner.Name + ":" + f.Config.GitConfig.BranchName,
			Base:     "master",
			Body:     f.Config.PullRequestConfig.Body,
			BodyFile: "",
			OrgFork:  f.Config.ForkConfig.Org,
		},
		// Number: 0,
		// ID:     0,
	}
	f.Logger.I("PRing", r.String())
	err = f.Repoer.CreatePullRequest(pr, f.Config.GitConfig.BranchNameBase) // pass BranchNameBase to check for pre-existing PRs
	if err != nil {
		return
	}
	outcome.PROutcome = pr
	outcome.Status = remote.PullRequested
	f.prsTally++
}
