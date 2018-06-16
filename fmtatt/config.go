package fmtatt

import (
	"github.com/rotblauer/gofmt-att/fmter"
	"github.com/rotblauer/gofmt-att/git"
	"github.com/rotblauer/gofmt-att/identity"
	"github.com/rotblauer/gofmt-att/logger"
	"github.com/rotblauer/gofmt-att/persist"
	"github.com/rotblauer/gofmt-att/remote"
	"github.com/rotblauer/gofmt-att/walk"
)

type Config struct {
	Identity          *identity.AuthIdentity // eg. whilei, ETCDEVTeam, etc.
	Pacing            Pacing
	RepoProvider      string
	ReposSpec         *remote.ListSpec
	Fmters            []*fmter.FmtConfig
	GitConfig         *git.GitConfig
	ForkConfig        *remote.ForkConfig
	PullRequestConfig *remote.PullRequestConfig
	PersistConfig     []*persist.PersistenceConfig
	WalkPattern       *walk.WalkPattern
	Logs              []*logger.LogConfig
}

type Pacing struct {
	MaxPRs                 int
	MininumPRSpreadMinutes int64 // minutes
}
