package common

import (
	"github.com/rotblauer/gofmt-att/fmtatt"
)

var DefaultFmtAttConfig = fmtatt.Config{
	RepoProvider: "Github",
	Identity: fmtatt.AuthIdentity{
		Username: "whilei",
		RawToken: "",
		EnvToken: "",
	},
	Fmt:               []fmtatt.FmtConfig{fmtatt.DefaultFmter},
	ReposFilter:       fmtatt.DefaultReposSpec,
	PullRequestConfig: fmtatt.DefaultPRConfig,
	WalkPattern:       fmtatt.DefaultWalkPattern,
	PersistConfig:     []fmtatt.PersistenceConfig{fmtatt.DefaultPersistenceConfig},
	Logs:              []fmtatt.LogConfig{fmtatt.DefaultLogConfig},
}
