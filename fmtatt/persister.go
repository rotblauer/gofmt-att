package fmtatt

import (
	"github.com/rotblauer/gofmt-att/remote"
	"github.com/rotblauer/gofmt-att/persist"
)

func (f *FmtAtt) mustGetState() persist.PersistentState {
	state, err := f.Persister.GetStateLeafs()
	if err != nil {
		if _, ok := err.(persist.ErrKeyNotFound); !ok {
			f.Logger.F(err)
		}
		gen := f.Config.WalkPattern.Genesis
		err := f.Persister.SetGenesis(gen)
		if err != nil {
			f.Logger.F(err)
		}
		state, err = f.Persister.GetStateLeafs()
		if err != nil {
			f.Logger.F(err)
		}
	}
	return state
}

func (f *FmtAtt) mustGetRepoOutcome(r *remote.RepoT) (o *remote.Outcome) {
	o, err := f.Persister.GetRepoOutcome(r)
	if err != nil {
		return nil // dangerzone
		// if _, ok := err.(persist.ErrKeyNotFound); !ok {
		// 	f.Logger.F(err)
		// }
	}
	return
}

func (f *FmtAtt) mustPutRepoOutcome(r *remote.RepoT, o *remote.Outcome) {
	if err := f.Persister.PutRepoOutcome(r, o); err != nil {
		f.Logger.F(err)
	}
}