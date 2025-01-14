package ethtrie

import (
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/v2/commitment/iavl"
	"cosmossdk.io/store/v2/db"
	"path/filepath"
)

var (
	StoreKeyUser  = "user-sk"
	StoreKeyOrder = "order-sk"
)

func newIVAL(dir string, disk bool) *iavl.IavlTree {
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
			filepath.Join(dir, "data"),
			nil,
		)
		if err != nil {
			panic(err)
		}
	}

	pdb := db.NewPrefixDB(rawdb, []byte(StoreKeyUser))
	opt := &iavl.Config{
		CacheSize:              500_000,
		SkipFastStorageUpgrade: true,
	}
	tree := iavl.NewIavlTree(pdb, log.NewNopLogger(), opt)
	return tree
}
