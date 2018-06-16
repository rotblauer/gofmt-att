package persist
//
// import (
// 	"github.com/dgraph-io/badger"
// 	"github.com/rotblauer/gofmt-att/remote"
// 	"encoding/json"
// )
//
// func (b BadgerPersistence) migrateFuckingLeaves() {
// 	err := b.db.Update(func(tx *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		opts.PrefetchSize = 10
// 		it := tx.NewIterator(opts)
// 		defer it.Close()
// 		for it.Rewind(); it.Valid(); it.Next() {
// 			item := it.Item()
// 			k := item.Key()
// 			v, err := item.Value()
// 			if err != nil {
// 				return err
// 			}
//
// 			var o *remote.Outcome
// 			err = json.Unmarshal(v, &o)
// 			if err != nil {
// 				return err
// 			}
//
//
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return
// }