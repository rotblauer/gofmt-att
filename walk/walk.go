package walk

import (
	"github.com/rotblauer/gofmt-att/persist"
	"github.com/rotblauer/gofmt-att/remote"
	"math/rand"
	"time"
	"sync"
	"errors"
	"fmt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// WalkProvider implements the logic required to select a next batch of repos to fetch.
type WalkProvider interface {
	// StepNext should be called when a list of repositories has been exhausted.
	// It should provide a node from which to fetch the next list of repositories.
	// In case a list of users is returned, that list should be traversed
	// name can be <(user|org)(|repo)>
	StepNext(pattern *WalkPattern, state persist.PersistentState, last Step, owners *remote.OwnerPool, repos *remote.RepoPool) (next Step)
}

type Step struct {
	Leaf remote.Leaf
	Err error
}

func (s Step) String() string {
	if s.Err != nil && s.Err.Error() != "" {
		return fmt.Sprintf("%s err=%s", s.Leaf.String(), s.Err.Error())
	}
	return fmt.Sprintf("%s", s.Leaf.String())
}

type WalkPattern struct {
	Name string
	Genesis           remote.Leaf
	QueryBranchConfig remote.QueryBranchConfig
}


func (wp *WalkPattern) NewWalker() DrunkenWalker {
	// establish a default
	drunkard := DrunkenWalker{
		lock: new(sync.RWMutex),
		WalkPattern: wp,
	}
	switch wp.Name {
	case "drunken":
		return drunkard
	}
	return drunkard
}

var (
	errIncompleteLeaf = errors.New("incomplete leaf")
	errFarFromHome = errors.New("far from home")
	errDeadEnd = errors.New("reached un-spannable gap; try raising your allowances for exploratoryness")
)



