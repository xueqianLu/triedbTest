package ethtrie

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/xueqianLu/triedbtest/testsuite"
	"math/big"
	"testing"
)

func getTrieDb(dir string, disk bool) ethdb.Database {
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

func TestHistoryTrie(t *testing.T) {
	dir := t.TempDir()
	db := getTrieDb(dir, true)
	defer db.Close()

	verifyUser := &types.StateAccount{
		CodeHash: []byte{0x1, 0x2, 0x3},
		Root:     common.HexToHash("0x123"),
		Nonce:    10018,
		Balance:  big.NewInt(111110),
	}
	verifyUserAddr := "122222222111111122"

	root := common.Hash{}
	snapSet := testsuite.NewSnapshotSet()
	for i := 0; i < 100; i++ {
		_, orderData := testsuite.GenerateAccount(200000)
		verifyUser.Balance = new(big.Int).Add(verifyUser.Balance, big.NewInt(int64(i)))
		orderData[verifyUserAddr] = testsuite.AccountData(verifyUser)
		tdb := trie.NewDatabase(db)
		// open tree, and set commit data to it.
		tree, err := trie.New(common.Hash{}, common.Hash{}, tdb)
		if err != nil {
			t.Fatalf("cannot create trie: %v", err)
		}
		for key, order := range orderData {
			if err := tree.TryUpdate([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				t.Fatalf("cannot update trie: %v", err)
			}
		}
		merged := trie.NewMergedNodeSet()
		newroot, nodes, err := tree.Commit(true)
		if err != nil {
			t.Fatalf("cannot commit trie: %v", err)
		}
		if err = merged.Merge(nodes); err != nil {
			t.Fatalf("cannot merge nodes: %v", err)
		}

		if err = tdb.Update(merged); err != nil {
			t.Fatalf("cannot update trie: %v", err)
		}
		if err = tdb.Commit(newroot, false, nil); err != nil {
			t.Fatalf("cannot commit trie: %v", err)
		}
		root = newroot
		snapSet.AddSnapshot(root, *verifyUser)
	}
	_, failed := snapSet.RangeSnapshot(func(sp testsuite.Snapshot) bool {
		vdb := trie.NewDatabase(db)
		tree, err := trie.New(common.Hash{}, sp.Root, vdb)
		if err != nil {
			t.Fatalf("cannot create trie: %v", err)
		}
		d, err := tree.TryGet([]byte(fmt.Sprintf("ux-%s", verifyUserAddr)))
		if err != nil {
			t.Fatalf("cannot get order: %v", err)
		}
		var order types.StateAccount
		if err := rlp.DecodeBytes(d, &order); err != nil {
			t.Fatalf("cannot decode order: %v", err)
		}
		if order.Balance.Cmp(sp.Account.Balance) != 0 {
			t.Fatalf("balance mismatch: have %v, want %v", order.Balance, sp.Account.Balance)
		}
		return true
	})
	if failed > 0 {
		t.Fatalf("failed to verify %d snapshots", failed)
	}
}

func BenchmarkTrieCommit(b *testing.B) {
	dir := b.TempDir()
	db := getTrieDb(dir, true)
	defer db.Close()

	root := common.Hash{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		_, orderData := testsuite.GenerateAccount(200000)
		b.StartTimer()
		tdb := trie.NewDatabase(db)
		// open tree, and set commit data to it.
		tree, err := trie.New(common.Hash{}, common.Hash{}, tdb)
		if err != nil {
			b.Fatalf("cannot create trie: %v", err)
		}
		for key, order := range orderData {
			if err := tree.TryUpdate([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				b.Fatalf("cannot update trie: %v", err)
			}
		}
		merged := trie.NewMergedNodeSet()
		newroot, nodes, err := tree.Commit(true)
		if err != nil {
			b.Fatalf("cannot commit trie: %v", err)
		}
		if err = merged.Merge(nodes); err != nil {
			b.Fatalf("cannot merge nodes: %v", err)
		}

		if err = tdb.Update(merged); err != nil {
			b.Fatalf("cannot update trie: %v", err)
		}
		if err = tdb.Commit(newroot, false, nil); err != nil {
			b.Fatalf("cannot commit trie: %v", err)
		}
		root = newroot
		b.StopTimer()
		size, _ := testsuite.GetDirSize(dir)
		b.Logf("dir size: %s", size.String())
		b.StartTimer()
	}
	_ = root
}
