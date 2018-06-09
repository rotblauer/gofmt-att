package fmtatt

import (
	"github.com/rotblauer/gofmt-att/fmtatt/providers/log"
	"log"
)

type Config struct {
	RepoProvider      string
	Identity          AuthIdentity // eg. whilei, ETCDEVTeam, etc.
	ReposFilter       RepoListSpec
	Fmt               []FmtConfig
	PullRequestConfig PullRequestConfig
	PersistConfig     []PersistenceConfig
	WalkPattern       WalkPattern
	Logs              []LogConfig
}

var DefaultFmtAttConfig = &Config{
	RepoProvider: "Github",
	Identity: AuthIdentity{
		Username: "whilei",
		RawToken: "",
		EnvToken: "",
	},
	Fmt:               []FmtConfig{DefaultFmter},
	ReposFilter:       DefaultReposSpec,
	PullRequestConfig: DefaultPRConfig,
	WalkPattern:       DefaultWalkPattern,
	PersistConfig:     []PersistenceConfig{DefaultPersistenceConfig},
	Logs:              []LogConfig{DefaultLogConfig},
}

type FmtAtt struct {
	Config *Config

	Repoer     RepoProvider
	Walker     WalkProvider
	Fmters     []Fmter
	Giter      GitProvider
	Persisters []PersistenceProvider
	Loggers    []Verbosably
}

func New(c *Config) *FmtAtt {
	log.Println("yep", c)

	f := FmtAtt{Config: c}

	switch c.RepoProvider {
	case "Github":
	default:
		log.Fatalln("unsupported repo provider type:", c.RepoProvider)
	}

	for _, l := range c.Logs {
		switch l.Logger {
		case "stderr":
			lo := logstderr.StdLogger{}
			lo.SetLevel(l.Level)
			f.Loggers = append(f.Loggers, lo)
		default:
			log.Fatalln("unsupported log type:", l)
		}
	}

	// TODO add them on as they get writ

	return nil
}

var (
	oneHalf  = float64(1) / float64(2)
	oneThird = float64(1) / float64(3)
)

