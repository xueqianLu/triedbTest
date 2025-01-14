package cosmos

import (
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/v2/commitment/iavl"
	"cosmossdk.io/store/v2/db"
)

var (
	StoreKeyUser  = "user-sk"
	StoreKeyOrder = "order-sk"
)

func NewRawDB(dir string, disk bool) (corestore.KVStoreWithBatch, error) {
	var (
		err   error
		rawdb corestore.KVStoreWithBatch
	)
	if !disk {
		rawdb = db.NewMemDB()
	} else {
		rawdb, err = db.NewDB(
			db.DBTypeGoLevelDB,
			"application",
			dir,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}
	return rawdb, nil
}

func NewIAVL(rawdb corestore.KVStoreWithBatch) *iavl.IavlTree {
	pdb := db.NewPrefixDB(rawdb, []byte(StoreKeyUser))
	opt := &iavl.Config{
		CacheSize:              200_000,
		SkipFastStorageUpgrade: true,
	}
	tree := iavl.NewIavlTree(pdb, log.NewNopLogger(), opt)
	return tree
}
