package ethtrie

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
)

func GetTrieDb(dir string, disk bool) ethdb.Database {
	var db ethdb.Database
	var err error
	if !disk {
		db = rawdb.NewMemoryDatabase()
	} else {
		db, err = rawdb.NewLevelDBDatabase(dir, 128, 1024, "", false)
		if err != nil {
			fmt.Printf("cannot create temporary database: %v", err)
		}
		return db
	}
	return rawdb.NewMemoryDatabase()
}
