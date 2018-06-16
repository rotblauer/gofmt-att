package remote

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
	"time"
	"fmt"
	"net/http"
	"errors"
	"github.com/kr/pretty"
)

const PerPage = 50 // 100 is max
const PageMax = 2

// GoogleGithubRepoProvider implements the Provider interface using the google/go-github API package.
type GoogleGithubRepoProvider struct {
	client     *github.Client
	httpC *http.Client
	ctx        context.Context
	username   string
	ownerSpecs *OwnerListSpec
	repoSpec   *RepoListSpec
	throttle   <-chan time.Time
	throttleD time.Duration
	startedAt time.Time
	reqCount  int
	quitting bool
}

func NewGoogleGithubProvider(username, token string) *GoogleGithubRepoProvider {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	ghc := github.NewClient(tc)

	startingRate := 1*time.Second
	gp := &GoogleGithubRepoProvider{
		client:    ghc,
		httpC: tc,
		ctx:       ctx,
		username:  username,
		throttle:  time.Tick(startingRate),
		throttleD: startingRate,
		startedAt: time.Now(),
		quitting:  false,
	}

	go func(gp *GoogleGithubRepoProvider) {
		// manage Github API rate limiting; max reqCount / hour MAX is 5000
		// calc max rate
		max := int64(5000)
		safeMaxPerMin := max / 10 / 60 // just to be safe...
		ticker := time.NewTicker(30*time.Second)
		for range ticker.C {
			fmt.Printf("(api requests=%d) repo provider heartbeat\n", gp.reqCount)
			if gp.quitting {
				fmt.Println("repoer quittting...")
				ticker.Stop()
				return
			}
			minutesRunning := time.Since(gp.startedAt).Minutes()
			if minutesRunning < 1 || gp.reqCount < 2 {
				continue
			}
			if rate := int64(gp.reqCount) / int64(minutesRunning); rate > safeMaxPerMin {
				gp.throttleD += 500*time.Millisecond
				gp.throttle = time.Tick(gp.throttleD)
			} else if rate < safeMaxPerMin/2 {
				gp.throttleD -= 500*time.Millisecond
				gp.throttle = time.Tick(gp.throttleD)
			}

			if gp.reqCount > 100 {
				panic("TOO MANY REQUESTS (DEVELOPMENT)")
			}
		}
	}(gp)

	return gp
}

func (rp *GoogleGithubRepoProvider) CastRepo(rawRepo interface{}) (repo *RepoT, ok bool) {
	r, ok := (rawRepo).(*github.Repository)
	if !ok {
		panic(pretty.Sprint(r))
		return nil, ok
	}
	repo = &RepoT{}
	repo.Name = r.GetName()
	if o, ok := rp.CastOwner(r.GetOwner()); ok {
		repo.Owner = o
	} else {
		return repo, false
	}
	repo.CloneUrl = r.GetCloneURL()
	repo.GitUrl = r.GetGitURL()
	repo.HTMLUrl = r.GetHTMLURL()
	return repo, true
}

func (rp *GoogleGithubRepoProvider) CastOwner(rawOwner interface{}) (owner *Owner, ok bool) {
	o, ok := (rawOwner).(*github.User)
	if ok {
		return &Owner{
			Name: o.GetLogin(),
			KindOf: o.GetType(),

		}, true
	}
	return nil, ok
}

func (rp *GoogleGithubRepoProvider) GetClient() *http.Client {
	return rp.httpC
}

func (rp *GoogleGithubRepoProvider) Cancel() {
	rp.quitting = true
}

func (rp *GoogleGithubRepoProvider) FetchRepo(r *RepoT) (rawRepo interface{}, err error) {
	rep, res, err := rp.client.Repositories.Get(rp.ctx, r.Owner.Name, r.Name)
	if ok, e := wrapGHRespErr(res, err); !ok {
		err = e
		return
	}
	rawRepo = rep
	return
}

func (rp *GoogleGithubRepoProvider) FetchLeaf(leaf Leaf, spec *RepoListSpec) (fetched []interface{}, err error) {
	if leaf.Header.IsListOfRepos() {
		return rp.getRepos(leaf, spec)
	}
	return rp.getUsers(leaf)
}

func (rp *GoogleGithubRepoProvider) getRepos(leaf Leaf, spec *RepoListSpec) (fetched []interface{}, err error) {
	var getter func(ctx context.Context, name string, options *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	var opts = &github.RepositoryListOptions{
		Visibility:  spec.Visibility,
		Affiliation: spec.Affiliation,
		Type:        "", // don't use with vis or aff
		Sort:        spec.SearchOptions.Sort,
		Direction:   spec.SearchOptions.Order,
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: PerPage,
		},
	}

	switch leaf.Header {
	case Starred:
		getter = func(ctx context.Context, name string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
			starred, res, err := rp.client.Activity.ListStarred(ctx, leaf.ID,
				&github.ActivityListStarredOptions{
					Sort:        opts.Sort,
					Direction:   opts.Direction,
					ListOptions: opts.ListOptions,
				})
			if err != nil {
				return nil, res, err
			}
			var repos []*github.Repository
			for _, s := range starred {
				repos = append(repos, s.Repository)
			}
			return repos, res, err
		}
	case Authored:
		getter = rp.client.Repositories.List
	case OrgRepos:
		getter = func(ctx context.Context, name string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
			return rp.client.Repositories.ListByOrg(ctx, name, &github.RepositoryListByOrgOptions{
				// Type:        "",
				ListOptions: opts.ListOptions,
			})
		}
	default:
		panic("TODO you shouldn't panic even if you think the code is unreachable" + leaf.Header)
	}
	for {
		rp.reqCount++
		repos, resp, gherr := getter(rp.ctx, leaf.ID, opts)
		if ok, e := wrapGHRespErr(resp, gherr); !ok {
			err = e
			return
		}
		fmt.Println("page", opts.Page, "got", len(repos), "repos")
		for _, r := range repos {
			fetched = append(fetched, r)
		}
		if resp.NextPage == 0 || resp.NextPage > PageMax {
			return
		}
		opts.Page = resp.NextPage
	}
	return
}

func (rp *GoogleGithubRepoProvider) getUsers(leaf Leaf) (fetched []interface{}, err error) {
	var getter func(context.Context, string, *github.ListOptions) ([]*github.User, *github.Response, error)
	var opts = &github.ListOptions{
		Page:    0,
		PerPage: PerPage,
	}
	switch leaf.Header {
	case Followers:
		getter = rp.client.Users.ListFollowers
	case Following:
		getter = rp.client.Users.ListFollowing
	case Stargazers:
		getter = func(ctx context.Context, name string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
			s := strings.Split(name, "/")
			stargazers, res, err := rp.client.Activity.ListStargazers(rp.ctx, s[0], s[1], opts)
			if stargazers == nil || err != nil {
				return nil, res, err
			}
			var users []*github.User
			for _, s := range stargazers {
				users = append(users, s.User)
			}
			return users, res, err
		}
		// // requires push authorization
	// case Collaborators:
	// 	getter = func(ctx context.Context, name string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
	// 		s := strings.Split(name, "/")
	// 		lopts := &github.ListCollaboratorsOptions{
	// 			Affiliation: "",
	// 			ListOptions: *opts,
	// 		}
	// 		return rp.client.Repositories.ListCollaborators(rp.ctx, s[0], s[1], lopts)
	// 	}
	case Members:
		getter = func(ctx context.Context, name string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
			lopts := &github.ListMembersOptions{
				PublicOnly:  false,
				Filter:      "",
				Role:        "",
				ListOptions: *opts,
			}
			return rp.client.Organizations.ListMembers(rp.ctx, name, lopts)
		}
	default:
		panic("TODO you shouldn't panic even if you think the code is unreachable" + leaf.Header)
	}
	for {
		rp.reqCount++
		users, resp, gherr := getter(rp.ctx, leaf.ID, opts)
		if ok, e := wrapGHRespErr(resp, gherr); !ok {
			err = e
			return
		}
		for _, u := range users {
			fetched = append(fetched, u)
		}
		if resp.NextPage == 0 || resp.NextPage > PageMax {
			break
		}
		opts.Page = resp.NextPage
	}
	return
}

func (rp *GoogleGithubRepoProvider) ForkRepo(config *ForkConfig, oR *RepoT) (origAndForked [2]*RepoT, err error) {
	origAndForked[0] = oR

	// <-rp.throttle
	rp.reqCount++
	fmt.Println("actually forking")

	// if the user belongs to the organization owning this repo, we'll still fork it if possible
	// if user has already forked the repo
	// if repo _belongs_ to auth'd user
	// GOTCHA we should only fork if this isn't the owner's repo...
	r, res, e := rp.client.Repositories.Get(rp.ctx, rp.username, oR.Name)
	if ok, _ := wrapGHRespErr(res, e); ok {
		if r != nil && r.GetGitURL() != "" {
			fmt.Println("not forking; ", rp.username+"/"+oR.Name, "exists")
			fmt.Println("--> ", r.GetHTMLURL())
			forkR, ok := rp.CastRepo(r)
			if !ok {
				panic("could not cast got repo")
			}
			origAndForked[1] = forkR
			return
		}
	}

	fmt.Println("no user fork/owned repo already, creating fork")
	rp.reqCount++
	r, res, err = rp.client.Repositories.CreateFork(rp.ctx, oR.Owner.Name, oR.Name, &github.RepositoryCreateForkOptions{
		Organization: config.Org,
	})

	// This method might return an *AcceptedError and a status code of 202.
	// This is because this is the status that GitHub returns to signify that it is now
	// computing creating the fork in a background task.
	// In this event, the Repository value will be returned, which includes the details
	// about the pending fork. A follow up request, after a delay of a second or so,
	// should result in a successful request.
	_, schedok := err.(*github.AcceptedError)
	for schedok {
		fmt.Println("sleeping for fork to happen")
		// Forking a Repository happens asynchronously. Therefore, you may have to wait a short period before accessing the git objects. If this takes longer than 5 minutes, be sure to contact GitHub support.
		time.Sleep(30*time.Second)
		own := rp.username
		if config.Org != "" {
			own = config.Org
		}
		fmt.Println("checking for user repo")
		rp.reqCount++
		_, res, err = rp.client.Repositories.Get(rp.ctx, own, oR.Name)
		_, schedok = err.(*github.AcceptedError)
	}
	if ok, e := wrapGHRespErr(res, err); !ok {
		err = e
		return
	}

	forkR, ok := rp.CastRepo(r)
	if !ok {
		panic("could not cast repo")
	}

	origAndForked[1] = forkR
	return
}

func (rp *GoogleGithubRepoProvider) CreatePullRequest(pr *PullRequestT, branchNameBase string) (err error) {
	rp.reqCount++
	// check for open prs matching what we're about to make.
	// NO DUPLICATE PRs!
	prs, res, err := rp.client.PullRequests.List(rp.ctx, pr.Owner.Name, pr.Name, &github.PullRequestListOptions{ListOptions: github.ListOptions{PerPage:PageMax}}) // State:"open" // either, really, at this point
	if ok, e := wrapGHRespErr(res, err); !ok {
		if len(prs) != 0 {
			err = e
			return
		}
	}

	prExists := func(p *github.PullRequest, pr *PullRequestT) bool {
		head := p.GetHead()
		cond := head.GetUser().GetLogin() == rp.username
		cond = cond && p.GetHead().GetRef() != ""
		cond = cond && strings.Contains(head.GetRef(), branchNameBase)
		cond = cond || p.GetTitle() == pr.Title
		cond = cond || p.GetBody() == pr.Body
		return cond
	}

	for _, p := range prs {
		if prExists(p, pr) {
			return errors.New("already PRed -> " + p.GetHTMLURL())
		}
	}

	rp.reqCount++
	fmt.Println("actually PRing")
	githubPR, res, err := rp.client.PullRequests.Create(rp.ctx, pr.Owner.Name, pr.Name, &github.NewPullRequest{
		Title:               github.String(pr.Title),
		Head:                github.String(pr.Head), // eg. rotblauer:fmt-att
		Base:                github.String("master"),
		Body:                github.String(pr.Body), // if was a body file, it will have been read and the value assigned to this field
		Issue:               nil,
		MaintainerCanModify: github.Bool(true),
	})
	if ok, e := wrapGHRespErr(res, err); !ok {
		err = e
		return
	}
	pr.Number = githubPR.GetNumber()
	pr.ID = githubPR.GetID()
	return


}
