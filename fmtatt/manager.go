package fmtatt

import (
	"github.com/rotblauer/gofmt-att/fmter"
	"github.com/rotblauer/gofmt-att/git"
	"github.com/rotblauer/gofmt-att/logger"
	"github.com/rotblauer/gofmt-att/persist"
	"github.com/rotblauer/gofmt-att/remote"
	"github.com/rotblauer/gofmt-att/walk"
	"os"
	"os/signal"
	"syscall"
	"time"
	"sync/atomic"
	"github.com/kr/pretty"
	"runtime"
)

/*
TODO: use an event mux to get events from stuff
*/

const (
	// Max number of owners to keep in mem at once.
	// Not really a big deal, just safe to keep a cap on it.
	ownerPoolMax = 100
)

var (
	// TODO assign these with optional CLI flags
	// Start actively trying to get more repos in queue in case at or below.
	repoQueueLowWater = 5 // runtime.GOMAXPROCS(0)
	// Clone or be cloning no more than this many repos at a time.
	cloningWorkLimit = runtime.GOMAXPROCS(0)

	// PTAL FIXME WTF
	cBufferSize = 100*(remote.PerPage * remote.PageMax) + repoQueueLowWater
)

type DryRunT bool

var DryRun DryRunT = true
var WetRun DryRunT = false

type FmtAtt struct {
	Config *Config

	Repoer remote.Provider
	Walker walk.WalkProvider
	Fmters []*fmter.FmtConfig
	Giter  git.GitProvider

	// This is kind of a strange idea. See interfacer.go
	// Interfacer *FmtAttInterfacer
	// - Persisters []persist.PersistenceProvider
	// - Loggers    []log.Verbosably

	// For now just keepin it simple.
	Persister persist.PersistenceProvider
	Logger    logger.Verbosably

	doFetchChan chan persist.PersistentState
	striperChan chan *remote.RepoT

	workerChan chan *remote.RepoT
	contributorChan chan struct{
		r *remote.RepoT
		o *remote.Outcome
	}

	quit          chan struct{}
	quitting bool
	pause         bool

	prIntervalMin time.Duration
	prsTally int

	persistentStateChanger func(st *persist.PersistentState, l remote.Leaf)

	repoPool    *remote.RepoPool
	ownerPool   *remote.OwnerPool
	workingPool *workPool

	dryRun [3]DryRunT // fork, push, pr
}

func (f *FmtAtt) Go(dryRun [3]DryRunT) {
	f.dryRun = dryRun

	// set defaults
	f.Logger.SetLevel(f.Config.Logs[0].Level) // yuck, cuz the multi-interfaces thing

	go func() {
		var sigc = make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigc)
		sig := <-sigc
		f.Logger.Ef("%s: shutting down", sig.String())
		f.quit <- struct{}{}
	}()

	pause := func(d time.Duration) {
		f.pause = true
		time.Sleep(d)
		f.pause = false
	}

	printStatus := func(state persist.PersistentState) {
		f.Logger.If("| PRs: %d", f.prsTally)
		f.Logger.If("| pool:working=%d %v", f.workingPool.len(), f.workingPool)
		f.Logger.If("| pool:owners=%d", f.ownerPool.Len())
		f.Logger.If("| pool:repos=%d", f.repoPool.Len())
		f.Logger.If("| %s", state.String())
	}

	s := f.mustGetState()
	go f.workerLoop()
	f.loadResources()

	printStatus(s)

	ticker := time.Tick(30 * time.Second)
	lastCheckedPRTally := 0
	lastSleep := time.Now()

	var fetching int32
	for {

		// water
		if l := f.repoPool.Len(); l < repoQueueLowWater && atomic.LoadInt32(&fetching) == 0 {
			// gotta query
			// get state
			state, err := f.Persister.GetStateLeafs()
			if err != nil {
				panic(err)
			}

			// replenish our owner pool if necessary
			if f.ownerPool.Len() <= 1 {
				os, err := f.Persister.GetOwners()
				if err != nil {
					f.Logger.F(err)
				}
				for i := 0; i < 10; i++ {
					f.ownerPool.Push(os[i])
				}
			}

			step := f.Walker.StepNext(f.Config.WalkPattern, state, walk.Step{Leaf: state.Current}, f.ownerPool, f.repoPool)
			if step.Err != nil {
				f.Logger.E("step err:", step.Err)
				continue
			}

			// get next step
			state, err = f.Persister.PutCurrentLeaf(step.Leaf, f.persistentStateChanger)
			if err != nil {
				f.Logger.F(err)
			}
			f.doFetchChan <- state
		}

		select {
		case state := <-f.doFetchChan:
			atomic.StoreInt32(&fetching, 1)
			f.fetch(state)
			//  f.workerChan <- [filtered repos]
			atomic.StoreInt32(&fetching, 0)

		case struc := <- f.contributorChan:
			f.Logger.I("contributor channel received %s %s", struc.r.String(), struc.o.String())
			go f.contributorFn(struc.r, struc.o)

		case wantRepo := <- f.striperChan:
			raw, err := f.Repoer.FetchRepo(wantRepo) // get github response for owner/name ref
			if err != nil {
				f.Logger.E("ERR FETCHING REPO:", err.Error())
			}
			rr, ok := f.Repoer.CastRepo(raw)
			if !ok {
				f.Logger.E("COULD NOT CAST:", pretty.Sprint(raw), "\nFROM:", wantRepo.String())
				continue
			}
			ok, outcome, filterErr := f.filterRawRepo(rr, raw)
			f.mustPutRepoOutcome(rr, outcome)

			if filterErr != nil {
				f.Logger.D("~ striper: ", rr.String(), outcome.String())
				continue
			}
			if !ok {
				continue
			}
			// set repo AND owner to pool
			f.Persister.PutOwner(rr.Owner)
			f.repoPool.Set(rr, outcome)
			if !f.ownerPool.Has(rr.Owner) {
				f.ownerPool.Push(rr.Owner)
			}
			f.workerChan <- rr // fire off to worker chan
			f.Logger.If("+ %s", rr.String())

		case <-ticker:
			// max out pr allowance
			if f.Config.Pacing.MaxPRs > 0 && f.prsTally > f.Config.Pacing.MaxPRs {
				f.quit<-struct{}{}
				return
			}
			// pause for pr minimum interval
			// noting that it's a ROUGH interval; max working pool allowance prs may still succeed
			if f.prsTally > lastCheckedPRTally && (f.prIntervalMin > 0 && time.Since(lastSleep) > f.prIntervalMin){
				f.Logger.W("%d PRs, sleeping %v", f.prIntervalMin.Round(time.Second))
				f.Logger.W("be advised; goroutines will be allowed to finish their tasks")
				pause(f.prIntervalMin)
				lastCheckedPRTally = f.prsTally
				lastSleep = time.Now()
			}
			state, err := f.Persister.GetStateLeafs()
			if err != nil {
				f.Logger.F(err)
			}
			printStatus(state)

		case <-f.quit:
			f.Repoer.Cancel()
			p := f.Config.GitConfig.BasePath
			if p != "" && p != "/" && p != "." {
				if err := os.RemoveAll(p); err != nil {
					panic(err)
				}
			}
			f.quitting = true
			return

		default:
		}
	}
}


func (f *FmtAtt) loadResources() {
	startingState := f.mustGetState() // get last or init from genesis
	f.Logger.I(startingState.String())

	o := startingState.Current.GetOwner()
	if startingState.Current.Header.BranchesFromOrg() {
		o.KindOf = "Organization"
	} else {
		o.KindOf = "User"
	}
	f.ownerPool.Push(o)

	// set up repos
	rs, e := f.Persister.GetRepos(func(outcome *remote.Outcome) bool {
		outcomeExpired := outcome.Timestamp.Before(time.Now().Add(-f.Config.ReposSpec.FmtExpiration))

		outcomeUnfinished := outcome.Status > remote.Clean // either cleared for fmting, with no progress otherwise
		outcomeUnfinished = outcomeUnfinished && outcome.Status < remote.PullRequested
		outcomeUnfinished = outcomeUnfinished && outcome.Error == "" // b/c it should be

		return outcome.Status == remote.Valid || outcomeExpired || outcomeUnfinished
	})
	if e != nil {
		f.Logger.F(e)
	}

	for _, r := range rs {
		outcome := f.mustGetRepoOutcome(r) // FIXME
		// we have to do this to get the owner type
		// FIXME this is dumb
		o, err := f.Persister.GetOwner(r.Owner.Name)
		if err != nil {
			f.Logger.F(err)
		} else if o == nil {
			// FIXME
			f.Logger.F("could not get owner from repo", r.Owner)
			continue
		}
		r.Owner = o
		f.repoPool.Set(r, outcome)
		if !f.ownerPool.Has(o) && o.KindOf != "" {
			f.ownerPool.Push(o)
		}
		f.workerChan <- r
	}
	return
}

