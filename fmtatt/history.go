package fmtatt

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
