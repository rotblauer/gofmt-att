package remote

import (
	"time"
	"net/http"
)

type Provider interface {
	GetClient() (*http.Client)
	FetchRepo(repo *RepoT) (rawRepo interface{}, err error)
	FetchLeaf(leaf Leaf, spec *RepoListSpec) (fetched []interface{}, err error)
	FilterRepo(rawRepo interface{}, spec *RepoListSpec) (err *ErrFilteredT)
	FilterOwner(rawOwner interface{}, spec *OwnerListSpec) (err *ErrFilteredT)
	ToRepo(rawRepo interface{}) (repo *RepoT, ok bool)
	ToOwner(rawUser interface{}) (owner *Owner, ok bool)
	ForkRepo(config *ForkConfig, rs *RepoT) (repos [2]*RepoT, err error)
	CreatePullRequest(pr *PullRequestT, branchNameBase string) (err error)
	Cancel()
}

type Resource struct {
	Owner *Owner
	Repo  *RepoT
}

type ListSpec struct {
	*RepoListSpec
	*OwnerListSpec
}

type MatchTextSpec struct {
	WhiteList []string
	BlackList []string
}

type MatchNSpec struct {
	Min, Max int64
}

type MatchTimeSpec struct {
	Earliest, Latest time.Time
}
