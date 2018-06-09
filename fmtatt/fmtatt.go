package fmtatt

import (
	"log"
)

type Config struct {
	RepoProvider      string
	Identity          AuthIdentity // eg. whilei, ETCDEVTeam, etc.
	ReposFilter       RepoListSpec
	Fmt               FmtConfig
	PullRequestConfig PullRequestConfig
	WalkPattern       WalkPattern
}

type FmtAtt struct {
	Config *Config

	Repoer RepoProvider
	Walker WalkProvider
	Fmters []Fmter
	Giter GitProvider
	Historyers []HistoryProvider

}

func New(c *Config) *FmtAtt {
	log.Println("yep", c)



	return nil
}