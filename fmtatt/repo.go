package fmtatt

import (
	"time"
)

type RepoProvider interface {
	// NewClient should be used to initialize client per package for global reuse there.
	NewClient(identity AuthIdentity) *interface{}
	GetRepos(reposFilter RepoListSpec) (repos []RepoT, err error)
	ForkRepo(rs RepoT) (repo RepoT, err error)
	CreatePullRequest(pr SimplePullRequestT) error
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
	RepoT: RepoT{
		Owner: "rotblauer",
		Name:  "gofmt-att",
	},
	Languages:  []string{"Go"},
	IsFork:     false,
	SortBy:     "updated",
	OrderBy:    "desc",
	Visibility: "visible",
}

type SimpleCommitConfig struct {
	Title      string
	Body       string
	AuthorName string
	Email      string
	Time       time.Time
}

type SimplePullRequestT struct {
	RepoT
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

type RepoT struct {
	Owner string
	Name  string
	CloneUrl string
	GitUrl   string
}

type RepoListSpec struct {
	RepoT
	Languages  []string
	IsFork     bool
	SortBy     string
	OrderBy    string
	Visibility string
}
