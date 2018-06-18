package fmtatt

import (
	"github.com/rotblauer/gofmt-att/fmter"
	"github.com/rotblauer/gofmt-att/git"
	"github.com/rotblauer/gofmt-att/identity"
	"github.com/rotblauer/gofmt-att/logger"
	"github.com/rotblauer/gofmt-att/persist"
	"github.com/rotblauer/gofmt-att/remote"
	"github.com/rotblauer/gofmt-att/walk"
	"math"
)

var DefaultFmtAttConfig = Config{
	RepoProvider: "Github",
	Pacing: Pacing{
		MaxPRs:                 50,
		MininumPRSpreadMinutes: 10,
	},
	Identity: &identity.AuthIdentity{
		Username: "whilei",
		RawToken: "",
		EnvToken: "GITHUB_TOKEN",
	},
	GitConfig: DefaultGitConfig,
	Fmters:    []*fmter.FmtConfig{DefaultFmter},
	ReposSpec: &remote.ListSpec{
		DefaultReposSpec,
		DefaultOwnersSpec,
	},
	ForkConfig:        DefaultForkConfig,
	PullRequestConfig: DefaultPRConfig,
	WalkPattern:       DefaultWalkPattern,
	PersistConfig:     []*persist.PersistenceConfig{DefaultPersistenceConfig},
	Logs:              []*logger.LogConfig{DefaultLogConfig},
}

var DefaultPersistenceConfig = &persist.PersistenceConfig{
	Name:     "badger",
	Endpoint: "/Users/ia/gofmt-att/badger.db",
}

var DefaultLogConfig = &logger.LogConfig{
	Level:  5, // 3 normally
	Logger: "stderr",
}

var DefaultForkConfig = &remote.ForkConfig{
	// if you want your forks to belong to an organization you have sufficient privileges for,
	// otherwise they will just be forked to your auth'd user.
	Org: "",
}

var DefaultPRConfig = &remote.PullRequestConfig{
	Title: "gofmt",
	Head:  "",
	Base:  "master",
	Body: `
Just ran

	gofmt -w .

on the project root. That's all.

> https://blog.golang.org/go-fmt-your-code
`,
	BodyFile: "",
	OrgFork: DefaultForkConfig.Org,
}

var DefaultFmter = &fmter.FmtConfig{
	Commands: []string{"gofmt -w"},
	PerFile: false, // the below lists are only in use when this is true
	Files: &fmter.FileList{
		WhiteList: []string{".go$"},
		BlackList: []string{""},
	},
	Dirs: &fmter.FileList{
		WhiteList: []string{""},
		BlackList: []string{"*vendor*", ".git"},
	},
}

var DefaultGitConfig = &git.GitConfig{
	Provider:        "gogit",
	BasePath:        "/Users/ia/gofmt-att/clones",
	GitCommitConfig: &git.GitCommitConfig{
		Title:          "all: gofmt", // this can be dynamically populated to be prefix with changes files or dirs, like core,stuff:
		BranchNameBase: "gofmt",
		Body:           `Run standard gofmt command on project root.

- go version go1.10.3 darwin/amd64`,
		AuthorName:     "ia",
		Email:          "isaac.ardis@gmail.com",
	},
	AddPaths: &remote.MatchTextSpec{
		WhiteList: []string{},
		BlackList: []string{`vendor\/`, `pkg\/`, "testdata", `assets\/`},
	},
	AddContent: &remote.MatchTextSpec{
		WhiteList: []string{},
		BlackList: []string{
			/*
			https://regex101.com/r/4ffW5H/2

// DO NOT EDIT!
// DONT EDIT!
// DON'T EDIT ME
// This file is auto-generated!
// DO NOT EDIT! This file is generated automatically
// DO NOT EDIT! This file is generated automaticaly
// This is a generated file, created automatically
// This file is generated automaticaly
			 */
			`^\/\/\s?(do(|\s?)n.?t\W?(edit|change|remove|add|tamper)|)(.*(auto(m|-)|)|generat|)`,



		},
	},
	StripeList: git.StripeList{
		`M\s*(|\/)vendor\/github\.com\/(?P<OWNER>\w*)\/(?P<REPO>[\w-]*\b)`,
	},
}

var DefaultReposSpec = &remote.RepoListSpec{
	Owner: &remote.MatchTextSpec{
		WhiteList: []string{},
		BlackList: []string{},
	},
	Name: &remote.MatchTextSpec{
		WhiteList: []string{},
		BlackList: []string{},
	},
	Description: &remote.MatchTextSpec{
		WhiteList: []string{},
		BlackList: []string{},
	},
	Language: &remote.MatchTextSpec{
		WhiteList: []string{"Go"},
		BlackList: []string{},
	},
	Conduct: &remote.MatchTextSpec{
		WhiteList: []string{},
		BlackList: []string{"NO ROBOTS", "ROBOTS NOT WELCOME"},
	},
	IsFork:   false,
	IsPrivate: false,
	Archived: false,
	StargazersCount: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	ForksCount: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	Size: &remote.MatchNSpec{
		Min: 0,
		Max: 1024 * 50, // kilobytes == ~50mb
	},
	NetworkCount: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	WatchersCount: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	Visibility:  "all",
	Affiliation: "owner,collaborator,organization_member",
	SearchOptions: &remote.SearchOptions{
		Sort:  "updated", // created, updated, pushed, full_name
		Order: "desc",    // asc or desc
	},
	FmtExpiry: 90, // days
}

// DefaultOwnersSpec applies for both users and orgs.
var DefaultOwnersSpec = &remote.OwnerListSpec{
	Owner: remote.Owner{
		Name: "",
	},
	FollowingN: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	FollowersN: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	PublicReposN: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	PublicGistsN: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
	CollaboratorsN: &remote.MatchNSpec{
		Min: 0,
		Max: math.MaxUint32,
	},
}

var DefaultWalkPattern = &walk.WalkPattern{
	Name: "drunken", // alternatives eventually might be 'methodical', 'clique'
	Genesis: remote.Leaf{
		Header: remote.OrgRepos, ID: "rotblauer",
	},
	QueryBranchConfig: remote.QueryBranchConfig{

		// How far the stepper should be allowed to walk from genesis.
		// Measured in leaves fetched that provide a list of repos (as opposed to a list of owners).
		MaxDistance: 60000,

		// Distribute branch leaf weights evenly
		Branching: func() map[remote.LeafHeader]float64 {
			nLeaves := len(remote.AvailableLeafHeaders)
			n := float64(1) / float64(nLeaves)
			leaves := make(map[remote.LeafHeader]float64, nLeaves)
			for i := 0; i < nLeaves; i++ {
				leaves[remote.AvailableLeafHeaders[i]] = n
			}
			return leaves
		}(),
	},
}
