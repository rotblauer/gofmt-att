package remote

import (
	"fmt"
	"strings"
)

type LeafHeader string
type Leaf struct {
	Header LeafHeader
	ID string  // eg. rotblauer or rotblauger/go-fmtatt
}

func (l Leaf) String() string {
	return fmt.Sprintf("%s:%s", l.Header, l.ID)
}

func (l Leaf) GetOwner() *Owner {
	return &Owner{
		Name: strings.Split(l.ID, "/")[0],
	}
}

func (l Leaf) IsRepo() bool {
	return len(strings.Split(l.ID, "/")) == 2
}

type QueryBranchConfig struct {
	MaxDistance int // tricky
	Branching   map[LeafHeader]float64
}

// these are completely fucked. i don't understand what i'm trying to do here.
var Branches = map[LeafHeader][]LeafHeader{
	OrgRepos: {Members},
	Members: {Starred, Authored, Following, Followers},

	// Collaborators: {Starred, Authored, Following, Followers},
	Stargazers: {Starred, Authored, Following, Followers},

	Starred: {Stargazers},
	Authored: {Stargazers},
	Following: {Starred, Authored},
	Followers: {Starred, Authored},
}

func (l LeafHeader) IsListOfRepos() bool {
	switch l {
	case Starred, Authored, OrgRepos:
		return true
	}
	return false
}

func (l LeafHeader) BranchesFromOrg() bool {
	return l == OrgRepos || l == Members
}

func (l LeafHeader) BranchesFromUser() bool {
	return l == Starred ||
		l == Authored ||
		l == Following ||
		l == Followers
}

func (l LeafHeader) BranchesFromRepo() bool {
	return l == Stargazers
}

// Owner-user leafs
const (
	Followers LeafHeader = "user_followers"
	Following            = "user_following"
)

// Owner-repo leafs
const (
	Starred  LeafHeader = "user_starred"
	Authored            = "user_authored"
)

// Repo-user leafs
const (
	Stargazers    LeafHeader = "repo_stargazers"
	// Collaborators            = "repo_collaborators"
	// Contributors       = "repo.contributors"
)

// Org-user leafs
const (
	Members LeafHeader = "org_members"
)

// Org-repo leafs
const (
	OrgRepos LeafHeader = "org_owns"
)

// Owner-org leafs
// const (
// 	Memberships LeafHeader = "user.orgs"
// )

var AvailableLeafHeaders = []LeafHeader{
	Following, Followers,
	Starred, Authored,
	Stargazers, // Collaborators,
	Members,
	OrgRepos, // Memberships,
}
