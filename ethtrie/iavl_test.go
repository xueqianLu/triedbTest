package ethtrie

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
)

func TestIAVLCommit(t *testing.T) {
	dir := t.TempDir()

	verifyUser := &types.StateAccount{
		CodeHash: []byte{0x1, 0x2, 0x3},
		Root:     common.HexToHash("0x123"),
		Nonce:    10018,
		Balance:  big.NewInt(111110),
	}
	verifyUserAddr := "122222222111111122"

	snapSet := NewSnapshotSet()
	for i := 0; i < 100; i++ {
		db := newIVAL(dir, true)
		_, orderData := generateAccount(200000)
		verifyUser.Balance = new(big.Int).Add(verifyUser.Balance, big.NewInt(int64(i)))
		orderData[verifyUserAddr] = accountData(verifyUser)

		for key, order := range orderData {
			if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				t.Fatalf("cannot set iavl trie: %v", err)
			}
		}
		_, newversion, err := db.Commit()
		if err != nil {
			t.Fatalf("cannot commit iavl trie: %v", err)
		}
		snapSet.AddSnapshot(common.BytesToHash(big.NewInt(int64(newversion)).Bytes()), *verifyUser)
	}
	_, failed := snapSet.RangeSnapshot(func(sp snapshot) bool {
		vdb := newIVAL(dir, true)
		version := big.NewInt(0).SetBytes(sp.root.Bytes()).Int64()
		d, err := vdb.Get(uint64(version), []byte(fmt.Sprintf("ux-%s", verifyUserAddr)))
		if err != nil {
			t.Fatalf("cannot get order: %v", err)
		}
		var order types.StateAccount
		if err := rlp.DecodeBytes(d, &order); err != nil {
			t.Fatalf("cannot decode order: %v", err)
		}
		if order.Balance.Cmp(sp.account.Balance) != 0 {
			t.Fatalf("balance mismatch: have %v, want %v", order.Balance, sp.account.Balance)
		}
		return true
	})
	if failed > 0 {
		t.Fatalf("failed to verify %d snapshots", failed)
	}
}

func BenchmarkIAVLCommit(b *testing.B) {
	dir := b.TempDir()
	root := common.Hash{}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db := newIVAL(dir, true)
		_, orderData := generateAccount(200000)
		b.StartTimer()
		for key, order := range orderData {
			if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				b.Fatalf("cannot set key: %v", err)
			}
		}
		_, newversion, err := db.Commit()
		if err != nil {
			b.Fatalf("cannot commit: %v", err)
		}
		db.Close()
		_ = newversion

	}
	_ = root
}
