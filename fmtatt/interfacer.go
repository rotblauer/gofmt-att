package fmtatt

// import (
// 	"github.com/rotblauer/gofmt-att/remote"
// 	"github.com/rotblauer/gofmt-att/persist"
// 	"github.com/rotblauer/gofmt-att/log"
// 	"fmt"
// )
//
// /*
// These are some sketches around handling multiple interfaces through one manager.
// Or managing interfaces through one handler.
// Or manhandling faces.
// <-done
//  */
//
// type FmtAttInterfacer struct {
// 	Persisters []persist.PersistenceProvider // x
// 	Loggers    []log.Verbosably
// }
//
// // Just take first error (or value).
// func (f *FmtAttInterfacer) PutDidWalkOne(r remote.RepoT) (err error) {
// 	var c = make(chan error, len(f.Persisters))
// 	for _, x := range f.Persisters {
// 		go func() {
// 			c<-x.PutDidWalkOne(r)
// 		}()
// 	}
// 	for {
// 		if err == nil || c != nil {
// 			err = <-c
// 		}
// 		if c != nil {
// 			close(c)
// 		}
// 		break
// 	}
// 	return
// }
//
// // Prioritize data consistency
// func (f *FmtAttInterfacer) GetDidWalkOne(r remote.RepoT) (ok bool, err error) {
// 	var oks = make([]bool, 2)
// 	var errs = make([]error, 2)
// 	for i, x := range f.Persisters {
// 		oks[i], errs[i] = x.GetDidWalkOne(r)
// 	}
// 	if oks[0] != oks[1] || errs[0] != errs[1] {
// 		panic("data is not consistent: " + fmt.Sprintf("%v %v %v %v", oks[0], oks[1], errs[0], errs[1]))
// 	} else {
// 		return oks[0], errs[0]
// 	}
// }
//
// func (f *FmtAttInterfacer) PutDidFmtOne(r remote.PullRequestT) error {
// 	for _, x := range f.Persisters {
// 		PutDidFmtOne
// 	}
// }
//
// func (f *FmtAttInterfacer) GetDidFmtOne(r remote.PullRequestT) (pr *remote.PullRequestT, err error) {
// 	for _, x := range f.Persisters {
// 		GetDidFmtOne
// 	}
// }
//
// func (f *FmtAttInterfacer) PutDidPROne(r remote.PullRequestT) error {
// 	for _, x := range f.Persisters {
// 		PutDidPROne
// 	}
// }
//
// func (f *FmtAttInterfacer) GetDidPROne(r remote.PullRequestT) (pr *remote.PullRequestT, err error) {
// 	for _, x := range f.Persisters {
// 		GetDidPROne
// 	}
// }
//
// func (f *FmtAttInterfacer) PutDidLeaf(leaf remote.Leaf) error {
// 	for _, x := range f.Persisters {
// 		PutDidLeaf
// 	}
// }
//
// func (f *FmtAttInterfacer) GetDidLeaf(leaf remote.Leaf) (ok bool, err error) {
// 	for _, x := range f.Persisters {
// 		GetDidLeaf
// 	}
// }
//
// func (f *FmtAttInterfacer) SetGenesis(leaf remote.Leaf) error {
// 	for _, x := range f.Persisters {
// 		SetGenesis
// 	}
// }
//
// func (f *FmtAttInterfacer) PutCurrentLeaf(leaf remote.Leaf) (state persist.PersistentState, err error) {
// 	for _, x := range f.Persisters {
// 		PutCurrentLeaf
// 	}
// }
//
// func (f *FmtAttInterfacer) GetStateLeafs() (state persist.PersistentState, err error) {
// 	for _, x := range f.Persisters {
// 		GetStateLeafs
// 	}
// }
//
// func (f *FmtAttInterfacer) Close() {
// 	for _, x := range f.Persisters {
// 		Close
// 	}
// }


// // Two grammars for assigning multiple interfaces
// f.Interfacer = &FmtAttInterfacer{
// 	Persisters: func() (ps []persist.PersistenceProvider) {
// 		for _, p := range c.PersistConfig {
// 			switch p.Name {
// 			case "badger":
// 				ps = append(ps, persist.NewBadger(p))
// 			default:
// 				panic("FIXME")
// 			}
// 		}
// 		return
// 	}(),
// }
// for _, l := range c.Logs {
// 	switch l.Logger {
// 	case "stderr":
// 		lo := log.StdLogger{}
// 		lo.SetLevel(l.Level)
// 		f.Interfacer.Loggers = append(f.Interfacer.Loggers, lo)
// 	default:
// 		fmt.Println("unsupported log type:", l)
// 		os.Exit(1)
// 	}
// }