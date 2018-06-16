package fmtatt

import (
	"github.com/rotblauer/gofmt-att/persist"
	"github.com/rotblauer/gofmt-att/remote"
	"fmt"
	"os"
	"github.com/rotblauer/gofmt-att/git"
	"github.com/rotblauer/gofmt-att/logger"
	"path/filepath"
	"time"
	"io/ioutil"
)

func New(c *Config) *FmtAtt {
	fmt.Println("yep", c)

	f := &FmtAtt{
		Config: c,
		Fmters: c.Fmters,
		Walker: c.WalkPattern.NewWalker(),
		Repoer: func() remote.Provider {
			switch c.RepoProvider {
			case "Github":
				return remote.NewGoogleGithubProvider(c.Identity.Username, mustGetTokenFromConfig(c))
			default:
				fmt.Println("unsupported remote provider type:", c.RepoProvider)
				os.Exit(1)
			}

			return nil
		}(),
		Giter: func() git.GitProvider {
			switch c.GitConfig.Provider {
			case "gogit":
				return git.NewGoGit()
			default:
				panic("FIXME")
			}
			return git.NewGoGit()
		}(),
	}

	if l := c.Logs[0]; l.Logger != "stderr" {
		panic("GET YOUR SHIT TOGETHER LOGGER")
	} else {
		f.Logger = logger.StdLogger{}
	}

	if p := c.PersistConfig[0]; p.Name != "badger" {
		panic("WTF just use badger already")
	}

	// sanitize and ensure exists persistent endpoing, eg. db path
	persistEndpoint := f.Config.PersistConfig[0].Endpoint
	persistEndpoint = filepath.Clean(persistEndpoint)
	persistEndpoint, err := filepath.Abs(persistEndpoint)
	if err != nil {
		panic("ABSOLUTELY NO DB PATH")
	}
	if err := os.MkdirAll(filepath.Dir(persistEndpoint), os.ModePerm); err != nil {
		panic(`WHY DON'T YOU JUST GO HOME BALL`)
	}
	f.Config.PersistConfig[0].Endpoint = persistEndpoint
	f.Logger.Df("db path:%s", persistEndpoint)
	f.Persister = persist.NewBadger(f.Config.PersistConfig[0])

	// sanitize and ensure exists clone path
	gitBasePath := filepath.Clean(f.Config.GitConfig.BasePath)
	gitBasePath, err = filepath.Abs(gitBasePath)
	if err != nil {
		panic("CLEAN ABSOLUTE PATHS PEOPLE" + err.Error())
	}
	if err := os.MkdirAll(gitBasePath, os.ModePerm); err != nil {
		panic("THE CLONES ARE NOT COMING")
	}
	f.Config.GitConfig.BasePath = gitBasePath
	f.Logger.Df("clone path:%s", f.Config.GitConfig.BasePath)
	f.Giter.SetBase(f.Config.GitConfig.BasePath)

	// parse expiry time in days to time.Duration
	f.Config.ReposSpec.FmtExpiration = time.Duration(f.Config.ReposSpec.FmtExpiry) * time.Hour * 24

	// read pull request template file if spec'd
	if bf := f.Config.PullRequestConfig.BodyFile; bf != "" {
		bf = filepath.Clean(bf)
		bs, err := ioutil.ReadFile(bf)
		if err != nil {
			panic("no body file found: " + bf)
		}
		f.Config.PullRequestConfig.Body = string(bs)
	}

	// append git branch name with timestamp
	// this should sidestep any problems with branch name conflicts
	if f.Config.GitConfig.BranchNameBase == "" {
		panic("cant have empty branch name base")
	}
	f.Config.GitConfig.GitCommitConfig.BranchName = f.Config.GitConfig.GitCommitConfig.BranchNameBase + time.Now().Format("-2006-Jan-2-15-04") // Mon Jan 2 15:04:05 -0700 MST 2006
	if f.Config.GitConfig.GitCommitConfig.BranchName == "" {
		panic("empty branch name")
	}

	// set min spread between PRs
	f.prIntervalMin = time.Duration(f.Config.Pacing.MininumPRSpreadMinutes)*time.Minute

	// set git to use authy client from repoer
	// this was to resolve issues with getting the git client to push; wasn't work
	// currently just using hardcoded system CLI 'git push -u <remote> <branch>'
	// TODO. rm me or fix me
	// i'm actually rather keen to just leavae it as os/env 'git', since it makes
	// usage requirements simple: just make sure your environment can push without requiring a password
	// to the https:// git endpoints
	// TODO maybe someday enable ssh pushes
	// if err := f.Giter.SetClient(f.Repoer.GetClient()); err != nil {
	// 	panic(err)
	// }

	f.doFetchChan = make(chan persist.PersistentState, 1)
	f.striperChan = make(chan *remote.RepoT, cBufferSize)
	f.workerChan = make(chan *remote.RepoT, cBufferSize)
	f.contributorChan = make(chan struct{
		r *remote.RepoT
		o *remote.Outcome
	}, cBufferSize)

	f.quit = make(chan struct{}, 1)

	// callback for updating state when setting anew
	f.persistentStateChanger = func(st *persist.PersistentState, l remote.Leaf) {

		st.Last = st.Current
		st.Current = l

		if l.Header.IsListOfRepos() {
			st.Distance++
		}

		st.Steps++

		// real rough; not sure how to measure distance
		if st.Genesis.ID == l.ID {
			f.Logger.W("at genesis")
			// st.Distance = 0
		}
		f.Logger.If("state: %s", st.String())
	}

	f.repoPool = remote.NewRepoPool()
	f.ownerPool = &remote.OwnerPool{}
	f.workingPool = &workPool{}

	return f
}
