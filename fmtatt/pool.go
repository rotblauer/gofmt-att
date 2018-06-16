package fmtatt

import (
	"github.com/rotblauer/gofmt-att/remote"
)

type workPool []*remote.RepoT

// OwnerPool deals with the randomness from Walker
func (wp *workPool) push(r *remote.RepoT) {
	oop := *wp
	oop = append(oop, r)
	*wp = oop
}

func (wp *workPool) has(r *remote.RepoT) bool {
	wpp := *wp
	for _, v := range wpp {
		if v.Ref() == r.Ref() {
			return true
		}
	}
	return false
}

func (wp *workPool) splice(r *remote.RepoT) {
	oop := *wp
	index := 0
	for i, v := range oop {
		if v.Ref() == r.Ref() {
			index = i
			break
		}
	}
	if len(oop) <= 1 {
		oop = oop[:0]
	} else {
		oop = append(oop[:index], oop[index+1:]...)
	}
	*wp = oop
}

func (wp *workPool) pop() *remote.RepoT {
	opp := *wp
	o := opp[0]
	opp = opp[1:]
	*wp = opp
	return o
}

func (wp *workPool) len() int {
	oop := *wp
	l := len(oop)
	oop = nil
	return l
}