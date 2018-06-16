package walk

import (
	"sync"
	"github.com/rotblauer/gofmt-att/remote"
	"math/rand"
	"time"
	"github.com/rotblauer/gofmt-att/persist"
	"fmt"
	"reflect"
)

type DrunkenWalker struct {
	*WalkPattern
	lock *sync.RWMutex
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// StepNext should only use persist's getter functions.
// NOTE this depends on the manager delegating valid relationship between leaf and candidates.
func (wp DrunkenWalker) StepNext(pattern *WalkPattern, state persist.PersistentState, last Step, candidates *remote.OwnerPool, repoPool *remote.RepoPool) (next Step) {

	// seems kind of weird, but just a default
	next = Step{}
	next.Leaf.Header = last.Leaf.Header
	next.Leaf.ID = last.Leaf.ID


	returnToGenesis := func(e error) {
		fmt.Printf("--------- RETURNING TO GENESIS: %v -----------\n", e.Error())
		next.Leaf = state.Genesis
		next.Err = e
		return
	}

	// starting from genesis
	if state.Distance == 0 && state.Steps == 0 {
		return
	}

	if state.Distance >= wp.WalkPattern.QueryBranchConfig.MaxDistance {
		returnToGenesis(errFarFromHome)
		return
	}

	if candidates == nil {
		panic("NIL CANDIDATES")
	}

	roulette := 0
	for roulette <= 1000 {
		roulette++
		o := drunkenDrawOwner(candidates)
		if o == nil {
			returnToGenesis(errDeadEnd)
		} else if o.KindOf == "" {
			panic("empty user kind of" + o.Name)
		} else {
			next.Leaf.ID = o.Name // this can be overridden if leaf requires repo endpoint
		}

		h, ok := drunkenGetLeafWeighted(wp.WalkPattern, remote.Branches[last.Leaf.Header]...)
		if !ok {
			returnToGenesis(errDeadEnd)
			return
		}
		if o.KindOf == "Organization" && !h.BranchesFromOrg() {
			continue
		}
		if o.KindOf == "User" && !h.BranchesFromUser() {
			continue
		}
		if h.BranchesFromRepo() {
			if repoPool.Len() == 0 {
				continue
			} else if repoPool.Len() == 1 {
				next.Leaf.ID = repoPool.Random1().Ref()
				if !next.Leaf.IsRepo() {
					continue
				}
			}
		}
		if reflect.DeepEqual(last, next) {
			continue
		}
		next.Leaf.Header = h
		break
	}
	if roulette == 1000 {
		panic("NEXT TIME 2 INFINITY")
	}

	return next
}

func drunkenDrawOwner(os *remote.OwnerPool) (owner *remote.Owner) {
	owners := *os
	if len(owners) == 0 {
		return
	} else if len(owners) == 1 {
		return owners[0]
	}
	owner = owners[rand.Intn(len(owners))]
	owners = nil
	return
}

func drunkenGetLeafWeighted(pattern *WalkPattern, leafs ...remote.LeafHeader) (lh remote.LeafHeader, ok bool) {
	leafSlice := append([]remote.LeafHeader{}, leafs...)
	nLeaves := len(leafSlice)

	if nLeaves == 0 {
		panic("nleaves 0")
		return
	} else if nLeaves == 1 {
		lh = leafSlice[0]
		ok = true
		return
	}

	// get sum
	var sum float64
	for _, l := range leafSlice {
		sum += pattern.QueryBranchConfig.Branching[l]
	}
	if sum == 0 {
		panic("NO SUM NO GOOD")
		return
	}
	rat := 1/sum
	if rat < 1 {
		panic("baaaad config. all configs vals together should add to one, with caveats")
	}
	sum = 0
	rnd := rand.Float64()
	for i := 0; i < nLeaves; i++ {
		mul := pattern.QueryBranchConfig.Branching[leafSlice[i]]
		sum += rat * mul
		if rnd < sum {
			lh = leafSlice[i]
			break
		}
	}
	if lh == "" {
		panic("emptpy head leafer")
	}
	ok = true
	return
}
