package remote

import (
	"sync"
	"math/rand"
)

type OwnerPool []*Owner

// OwnerPool deals with the randomness from Walker
func (op *OwnerPool) Push(o *Owner) {
	oop := *op
	oop = append(oop, o)
	*op = oop
}

func (op *OwnerPool) Splice(o *Owner) (ow *Owner) {
	oop := *op
	index := 0
	for i, v := range oop {
		if v.Name == o.Name {
			index = i
			ow = v
			break
		}
	}
	if len(oop) <=1 {
		oop = oop[:0]
	} else {
		oop = append(oop[:index], oop[index+1:]...)
	}
	*op = oop
	return
}

func (op *OwnerPool) Has(o *Owner) bool {
	oop := *op
	for _, v := range oop {
		if v.String() == o.String() {
			return true
		}
	}
	oop = nil
	return false
}

func (op *OwnerPool) Pop() *Owner {
	opp := *op
	o := opp[0]
	opp = opp[1:]
	*op = opp
	return o
}

func (op *OwnerPool) Len() int {
	opp := *op
	i := len(opp)
	opp = nil
	return i
}


func (op *OwnerPool) AsStrings() (candidates []string) {
	oop := *op
	for _, o := range oop {
		candidates = append(candidates, o.String())
	}
	oop = nil
	return
}

type RepoPool struct {
	m map[*RepoT]*Outcome
	mut *sync.RWMutex
	lim int
}

func NewRepoPool() *RepoPool {
	return &RepoPool{
		m:   make(map[*RepoT]*Outcome),
		mut: new(sync.RWMutex),
		lim: 1000, // TODO
	}
}

// RepoPool is FIFO, for now
func (rp *RepoPool) Set(r *RepoT, o *Outcome) {
	rp.mut.Lock()
	rp.m[r] = o
	rp.mut.Unlock()
}

func (rp *RepoPool) Remove(r *RepoT) {
	rp.mut.Lock()
	delete(rp.m, r)
	rp.mut.Unlock()
}

func (rp *RepoPool) GetOutcome(r *RepoT) *Outcome {
	rp.mut.RLock()
	defer rp.mut.RUnlock()
	return rp.m[r]
}

func (rp *RepoPool) Len() int {
	rp.mut.RLock()
	defer rp.mut.RUnlock()
	return len(rp.m)
}

func (rp *RepoPool) Random1() *RepoT {
	rp.mut.RLock()
	defer rp.mut.RUnlock()
	l := len(rp.m)
	r := rand.Intn(l)
	i := 0
	var repo *RepoT
	for rr := range rp.m {
		if i == r {
			repo = rr
		}
		i++
	}
	return repo
}

func (rp *RepoPool) Repos() (repos []*RepoT) {
	rp.mut.RLock()
	defer rp.mut.RUnlock()
	for r := range rp.m {
		repos = append(repos, r)
	}
	return
}

func (rp *RepoPool) GetWhere(filter func(r *RepoT, o *Outcome) bool) (repos []*RepoT) {
	rp.mut.RLock()
	defer rp.mut.RUnlock()
	for rr, oo := range rp.m {
		if filter(rr, oo) {
			repos = append(repos, rr)
		}
	}
	return
}