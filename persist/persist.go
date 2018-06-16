package persist

import (
	"github.com/rotblauer/gofmt-att/remote"
	"fmt"
	"reflect"
)

type PersistenceConfig struct {
	Name     string // eg. bolt, kf, POST
	Endpoint string // eg. path/to/database or HTTP endpoint
}

type ErrKeyNotFound struct {
	error
}

type PersistentState struct {
	Genesis remote.Leaf
	Last remote.Leaf
	Current remote.Leaf
	Distance int
	Steps int64 // ha
}

func (ps PersistentState) IsGenesis() bool {
	return reflect.DeepEqual(ps.Genesis, ps.Current) && ps.Steps == 0
}

func (ps PersistentState) String() string {
	return fmt.Sprintf("(*)genesis=%s [last=%s --> current=%s] steps=%d", ps.Genesis, ps.Last, ps.Current, ps.Steps)
}

type PersistenceProvider interface {
	PutRepoOutcome(r *remote.RepoT, outcome *remote.Outcome) error
	GetRepoOutcome(r *remote.RepoT) (outcome *remote.Outcome, err error)

	PutOwner(o *remote.Owner) error
	GetOwner(name string) (owner *remote.Owner, err error)
	GetOwners() (owners []*remote.Owner, err error)

	GetRepos(withOutcome func(outcome *remote.Outcome) (matching bool)) (repos []*remote.RepoT, err error)

	SetGenesis(leaf remote.Leaf) error
	PutCurrentLeaf(leaf remote.Leaf, updater func(st *PersistentState, l remote.Leaf)) (state PersistentState, err error)
	GetStateLeafs() (state PersistentState, err error)

	Close()
}

type RetrospectiveProvider interface {
	GetOutcomes(outcome remote.Outcome) (repos []*remote.RepoT, err error)
	Close()
}