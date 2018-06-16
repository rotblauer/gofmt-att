package persist

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/rotblauer/gofmt-att/remote"
	"time"
	"strings"
	"log"
)

type BadgerPersistence struct {
	db *badger.DB
}

var (
	repoP = []byte("r~")
	ownP = []byte("o~")
	stateLeafsP = []byte("@~") // origin, last, current, steps
)

func NewBadger(c *PersistenceConfig) BadgerPersistence {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = c.Endpoint
	opts.ValueDir = c.Endpoint
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	bp := BadgerPersistence{
		db: db,
	}
	return bp
}

func repoToKey(r *remote.RepoT) []byte {
	if r.Owner.Name == "" {
		panic("empty owner")
	} else if r.Name == "" {
		panic("empty repo")
	}

	k := repoP
	k = append(k, []byte(r.Owner.Name)...)
	k = append(k, []byte("/")...)
	k = append(k, []byte(r.Name)...)

	return k
}

func keyToRepo(key []byte) *remote.RepoT {
	sor := string(key[len(repoP):])
	or := strings.Split(sor, "/")
	return &remote.RepoT{
		Owner:    &remote.Owner{
			Name:   or[0],
			KindOf: "",
		},
		Name:     or[1],
		Target:   "",
		CloneUrl: "",
		GitUrl:   "",
		HTMLUrl:  "",
	}
}

func ownerToKey(o *remote.Owner) []byte {
	if o.Name == "" {
		panic("empty owner")
	}

	k := ownP
	k = append(k, []byte(o.Name)...)
	return k
}

func keyToOwner(key []byte) *remote.Owner {
	return &remote.Owner{Name: string(key[len(ownP):])}
}

func (b BadgerPersistence) PutOwner(o *remote.Owner) (err error) {
	err = b.db.Update(func(tx *badger.Txn) error {
		return tx.Set(ownerToKey(o), []byte(o.KindOf)) // nil for now
	})
	return
}

func (b BadgerPersistence) GetOwners() (owners []*remote.Owner, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(ownP); it.ValidForPrefix(ownP); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
			o := keyToOwner(k)
			o.KindOf = string(v)
			owners = append(owners, o)
		}
		return nil
	})
	return
}

func (b BadgerPersistence) GetOwner(name string) (o *remote.Owner, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(ownP); it.ValidForPrefix(ownP); it.Next() {
			item := it.Item()
			k := item.Key()
			o = keyToOwner(k)
			if o.Name == name {
				v, err := item.Value()
				if err != nil {
					return err
				}
				o.KindOf = string(v)
				return nil
			}
		}
		return nil
	})
	return
}



func (b BadgerPersistence) PutRepoOutcome(r *remote.RepoT, outcome *remote.Outcome) (err error) {
	err = b.db.Update(func(tx *badger.Txn) error {
		outcome.Timestamp = time.Now()
		b, err := json.Marshal(outcome)
		if err != nil {
			return err
		}
		return tx.Set(repoToKey(r), b)
	})
	return
}

func (b BadgerPersistence) GetRepoOutcome(r *remote.RepoT) (outcome *remote.Outcome, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		v, err := tx.Get(repoToKey(r))
		if err == badger.ErrKeyNotFound {
			return ErrKeyNotFound{err}
		} else if err != nil {
			return err
		}
		b, err := v.Value()
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &outcome)
		if err != nil {
			return err
		}
		return err
	})
	return
}

func (b BadgerPersistence) GetRepos(withOutcome func(outcome *remote.Outcome) (matching bool)) (repos []*remote.RepoT, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(repoP); it.ValidForPrefix(repoP); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}

			var o = &remote.Outcome{}
			err = json.Unmarshal(v, o)
			if err != nil {
				log.Println(err)
				log.Fatalln(string(v))
				return err
			}

			if withOutcome(o) {
				// establish repo
				r := keyToRepo(k)
				repos = append(repos, r)
			}
		}
		return nil
	})
	return
}

func (b BadgerPersistence) SetGenesis(leaf remote.Leaf) (err error) {
	err = b.db.Update(func(tx *badger.Txn) error {
		state := PersistentState{
			Genesis: leaf,
			Last: leaf,
			Current: leaf,
			Steps: 0,
		}
		b, err := json.Marshal(&state)
		if err != nil {
			return err
		}
		return tx.Set(stateLeafsP, b)
	})
	return
}

func (b BadgerPersistence) PutCurrentLeaf(leaf remote.Leaf, changeState func(st *PersistentState, l remote.Leaf)) (state PersistentState, err error) {
	err = b.db.Update(func(tx *badger.Txn) error {
		item, err := tx.Get(stateLeafsP)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrKeyNotFound{err}
			}
			return err
		}
		b, err := item.Value()
		if err != nil {
			return err
		}
		var p = &PersistentState{}
		err = json.Unmarshal(b, p)
		if err != nil {
			return err
		}

		changeState(p, leaf)

		state = *p

		bb, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return tx.Set(stateLeafsP, bb)
	})
	return
}

func (b BadgerPersistence) GetStateLeafs() (state PersistentState, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(stateLeafsP)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrKeyNotFound{err}
			}
			return err
		}
		b, err := item.Value()
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &state)
		if err != nil {
			return err
		}
		return err
	})
	return
}

func (b BadgerPersistence) Close() {
	b.db.Close()
}
