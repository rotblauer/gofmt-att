package remote

import (
	"fmt"
	"strconv"
)

type PullRequestT struct {
	*RepoT
	*PullRequestConfig
	// on success, we store these
	Number    int
	ID        int64
}

type PullRequestConfig struct {
	Title    string
	Head     string `json:"-",toml:"-",yaml:"-"` // (branch name to Set Set changes to; create PR rom) Default to an automatic probably-unique one, like types-20180606
	Base     string `json:"-",toml:"-",yaml:"-"` // Default to 'master'
	Body     string
	BodyFile string
	OrgFork  string // Organization name
}

func (pr *PullRequestT) String() string {
	return fmt.Sprintf(`🎈︎ PR %s`, pr.RepoT.HTMLUrl + "/pulls/" + strconv.Itoa(pr.Number))
}