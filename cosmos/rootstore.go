package cosmos

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store/v2/db"
	"cosmossdk.io/store/v2/root"
	"path/filepath"
)

var (
	StoreKeyUser  = "user-sk"
	StoreKeyOrder = "order-sk"
)

func newStore(dir string) {
	opt := root.DefaultStoreOptions()
	opt.IavlConfig.CacheSize = 500_000
	rawdb, err := db.NewDB(
		db.DBTypeGoLevelDB,
		"application",
		filepath.Join(dir, "data"),
		nil,
	)
	if err != nil {
		panic(err)
	}
	cfg := &root.FactoryOptions{
		Logger:    log.NewNopLogger(),
		RootDir:   dir,
		Options:   opt,
		StoreKeys: []string{StoreKeyUser, StoreKeyOrder},
		SCRawDB:   rawdb,
	}
	store, err := root.CreateRootStore(cfg)
	if err != nil {
		panic(err)
	}
	_ = store

}
