package fmtatt

type PersistenceConfig struct {
	Name string // eg. bolt, kf, POST
	Endpoint string // eg. path/to/database or HTTP endpoint
}

var DefaultPersistenceConfig = PersistenceConfig{
	Name: "bolt",
	Endpoint: "/var/gofmt-att.db",
}

type PersistenceProvider interface {
	PutDidFmtOne(config Config, pr SimplePullRequestT) error
	GetDidFmtOne(config Config, repo RepoSpec) (pr *SimplePullRequestT, err error)

	PutDidFmtList(config Config, rs RepoListSpec) error
	GetDidFmtList(config Config) (rss []RepoListSpec, err error)

	PutDidWalkOne(config Config, r RepoSpec) error
	GetDidWalkOne(config Config, r RepoListSpec) (didWalk bool, err error)

	PutDidWalkList(config Config, rs RepoListSpec) error
	GetDidWalkList(config Config, rs RepoListSpec) (didWalk bool, err error)
}
