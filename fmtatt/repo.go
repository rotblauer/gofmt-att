package fmtatt

import (
	"context"
	"net/http"
	"time"
)

type RepoProvider interface {
	// NewClient should be used to initialize client per package for global reuse there.
	NewClient(*http.Client) interface{}
	GetRepos(ctx context.Context, reposFilter RepoListSpec) (repos []SimpleRemoteT, err error)
	ForkRepo(ctx context.Context, authID AuthIdentity) (repo SimpleRemoteT, err error)
	CreatePullRequest(ctx context.Context, pr SimplePullRequestT) error
}

var DefaultPRConfig = PullRequestConfig{
	Title: "gofmt: nice and clean. your's truly, the gofmt-att machine",
	Head:  "",
	Base:  "master",
	Body: `
Formattered with :heart: by [gofmt-att](https://github.com/rotblauer/gofmt-att).

> If we got it wrong, or there's a bug or something, please [let us know](https://github.com/rotblauer/gofmt-att/issues/new).
`,
	BodyFile: "",
}

var DefaultReposSpec = RepoListSpec{
	RepoSpec: RepoSpec{
		Owner: "rotblauer",
		Name:  "gofmt-att",
	},
	Languages:  []string{"Go"},
	IsFork:     false,
	SortBy:     "updated",
	OrderBy:    "desc",
	Visibility: "visible",
}

type SimpleRemoteT struct {
	CloneUrl string
	GitUrl   string
}

type SimpleCommitConfig struct {
	Title      string
	Body       string
	AuthorName string
	Email      string
	Time       time.Time
}

type SimplePullRequestT struct {
	AuthIdentity
	SimpleRemoteT
	Title string
	Head  string
	Base  string
	Body  string
	Diff  []byte
}

type PullRequestConfig struct {
	Title    string
	Head     string // Default to an automatic probably-unique one, like fmtatt-20180606
	Base     string // Default to 'master'
	Body     string
	BodyFile string
}

type RepoSpec struct {
	Owner string
	Name  string
	SimpleRemoteT
}

type RepoListSpec struct {
	RepoSpec
	Languages  []string
	IsFork     bool
	SortBy     string
	OrderBy    string
	Visibility string
}
