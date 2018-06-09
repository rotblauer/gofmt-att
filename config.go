package main

import (
	"net/http"
	"context"
	"time"
	"io"
)

type APIProviderName string

type AuthIdentity struct {
	Name string
	RawToken string
	EnvToken string
}

type RepoProvider interface {
	// NewClient should be used to initialize client per package for global reuse there.
	NewClient(*http.Client) interface{}
	GetRepos(ctx context.Context, reposFilter RepoListSpec) (repos []SimpleRemoteT, err error)
	ForkRepo(ctx context.Context, authID AuthIdentity) (repo SimpleRemoteT, err error)
	CreatePullRequest(ctx context.Context, pr SimplePullRequestT) error
}

type WalkProvider interface {
	StepNext(pattern WalkPattern, history HistoryProvider) (rs RepoListSpec, err error)
}

type fmtErrs struct {
	err error
	msg io.Reader
}
type fmted struct {
	msg io.Reader
}
type Fmter interface {
	Gofmt(fmtConfig FmtConfig, fmted chan fmted, errs chan fmtErrs) (done chan bool)
}

type GitProvider interface {
	Clone(remote SimpleRemoteT) (dirPath string, err error) // Path to repo (name), this will be used by the FmtConfig.Commands and FmtConfig.Target
	IsDirty(dirPath string) bool
	CreateBranch(dirPath string, branchName string) error
	AddRemote(dirPath string, remoteName string) error
	AddAndCommitAll(dirPath string, commit SimpleCommitConfig) (diff io.Reader, err error)
	PushAll(dirPath string, remote string) error
}

type HistoryProvider interface {
	PutDidFmtOne(config FmtAttConfig, pr SimplePullRequestT) error
	GetDidFmtOne(config FmtAttConfig, repo RepoSpec) (pr *SimplePullRequestT, err error)

	PutDidFmtList(config FmtAttConfig, rs RepoListSpec) error
	GetDidFmtList(config FmtAttConfig) (rss []RepoListSpec, err error)

	PutDidWalkOne(config FmtAttConfig, r RepoSpec) error
	GetDidWalkOne(config FmtAttConfig, r RepoListSpec) (didWalk bool, err error)

	PutDidWalkList(config FmtAttConfig, rs RepoListSpec) error
	GetDidWalkList(config FmtAttConfig, rs RepoListSpec) (didWalk bool, err error)
}

type FmtAttConfig struct {
	Provider          APIProviderName
	Identity          string // eg. whilei, ETCDEVTeam, etc.
	ReposFilter       RepoListSpec
	PullRequestConfig PullRequestConfig
	WalkPattern       WalkPattern
	Fmt               FmtConfig
}

type FmtConfig struct {
	Commands []string
	Target string
	PerFile bool
	FilesWhitelist []string
	FilesBlacklist []string
	DirsWhitelist []string
	DirsBlacklist []string
}

type SimpleRemoteT struct {
	CloneUrl string
	GitUrl string
}

type SimpleCommitConfig struct {
	Title string
	Body string
	AuthorName string
	Email string
	Time time.Time
}

type SimplePullRequestT struct {
	AuthIdentity
	SimpleRemoteT
	Title string
	Head string
	Base string
	Body string
	Diff []byte
}

type PullRequestConfig struct {
	Title string
	Head string // Default to an automatic probably-unique one, like fmtatt-20180606
	Base string // Default to 'master'
	Body string
	BodyFile string
}

type RepoSpec struct {
	Owner string
	Name string
	SimpleRemoteT
}

type RepoListSpec struct {
	RepoSpec
	Languages []string
	IsFork bool
	SortBy string
	Visibility string
}

type WalkPattern struct {
	WalkHumansPattern
	WalkReposPatterns
}

type WalkHumansPattern struct {
	Following bool
	Followers bool
	OrgMembers bool
}

type WalkReposPatterns struct {
	Starred bool
	Forked bool
	Authored bool
}

