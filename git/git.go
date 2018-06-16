package git

import (
	"github.com/rotblauer/gofmt-att/remote"
	"time"
	"net/http"
)

type GitProvider interface {
	SetBase(ppath string) (err error)
	SetClient(client *http.Client) (err error)
	Clone(repo *remote.RepoT) (err error) // Path to repo (name), this will be used by the FmtConfig.Commands and FmtConfig.Target
	Status(dirPath string) (dirty bool, gitStatus string, err error)
	Add(dirPath string, filePath string) (err error)
	CommitWithBranch(dirPath string, commit *GitCommitConfig) (hash, status string, err error)
	PushAll(dirPath, remoteUrl, branchName string) (err error)
}

type GitConfig struct {
	Provider string
	BasePath string
	*GitCommitConfig
	AddPaths *remote.MatchTextSpec
	StripeList StripeList
}

type GitCommitConfig struct {
	BranchNameBase string
	BranchName     string    `json:"-",toml:"-",yaml:"-"`
	Title          string
	Body           string
	AuthorName     string
	Email          string
	Time           time.Time `json:"-",toml:"-",yaml:"-"`
}

// TODO: StripeList
// A stripelist is a list of regex strings (MustCompile), that
// must yield named matches "OWNER" and "REPO".
// After fmting, all dirty files (white/blacklist agnostic) will be matched against this regex.
// Resulting OWNER/REPO results will be forwarded to the repo service for fetching.
// This lets you exclude changes against, say, vendor/ directories in your commit, and instead
// track down those sources for formatting themselves.
// eg.
// M\s*(|\/)vendor\/github\.com\/(?P<OWNER>\w*)\/(?P<REPO>[\w-]*\b)
// https://regex101.com/r/BweS5r/2
type StripeList []string