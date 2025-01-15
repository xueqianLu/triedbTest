package cosmos

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/triedbtest/testsuite"
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
	rawdb, err := NewRawDB(dir, true)
	if err != nil {
		logrus.WithError(err).Error("cannot create raw db")
		return
	}
	defer rawdb.Close()

	snapSet := testsuite.NewSnapshotSet()
	for i := 0; i < 100; i++ {
		db := NewIAVL(rawdb)
		latest, err := db.GetLatestVersion()
		if err != nil {
			t.Fatalf("cannot get latest version: %v", err)
		}
		if err := db.LoadVersion(latest); err != nil {
			t.Fatalf("cannot load version: %v", err)
		}
		_, orderData := testsuite.GenerateAccount(200000)
		verifyUser.Balance = new(big.Int).Add(verifyUser.Balance, big.NewInt(int64(i)))
		orderData[verifyUserAddr] = testsuite.AccountData(verifyUser)

		for key, order := range orderData {
			if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				t.Fatalf("cannot set iavl trie: %v", err)
			}
		}
		_, newversion, err := db.Commit()
		if err != nil {
			t.Fatalf("cannot commit iavl trie: %v", err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("cannot close iavl trie: %v", err)
		}
		snapSet.AddSnapshot(common.BytesToHash(big.NewInt(int64(newversion)).Bytes()), *verifyUser)
	}
	_, failed := snapSet.RangeSnapshot(func(sp testsuite.Snapshot) bool {
		vdb := NewIAVL(rawdb)
		version := big.NewInt(0).SetBytes(sp.Root.Bytes()).Int64()
		d, err := vdb.Get(uint64(version), []byte(fmt.Sprintf("ux-%s", verifyUserAddr)))
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
		vdb.Close()
		return true
	})
	if failed > 0 {
		t.Fatalf("failed to verify %d snapshots", failed)
	}
}

func BenchmarkIAVLCommit(b *testing.B) {
	dir := b.TempDir()
	rawdb, err := NewRawDB(dir, true)
	if err != nil {
		b.Fatalf("cannot create raw db: %v", err)
	}
	defer rawdb.Close()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db := NewIAVL(rawdb)
		_, orderData := testsuite.GenerateAccount(200000)
		latest, err := db.GetLatestVersion()
		if err != nil {
			logrus.WithError(err).Error("cannot get latest version")
			return
		}
		if err := db.LoadVersion(latest); err != nil {
			logrus.WithField("version", latest).WithError(err).Error("cannot load version")
			return
		}

		for key, order := range orderData {
			if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				b.Fatalf("cannot set key: %v", err)
			}
		}
		b.StartTimer()
		_, _, err = db.Commit()
		if err != nil {
			b.Fatalf("cannot commit: %v", err)
		}
		db.Close()
		b.StopTimer()
		size, _ := testsuite.GetDirSize(dir)
		b.Logf("dir size: %s", size.String())
		b.StartTimer()
	}
}

func BenchmarkIAVLCommitCustom(b *testing.B) {
	dir := b.TempDir()
	rawdb, err := NewRawDB(dir, true)
	if err != nil {
		b.Fatalf("cannot create raw db: %v", err)
	}
	defer rawdb.Close()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		db := NewIAVL(rawdb)
		_, orderData := testsuite.GenerateCustom(200000)
		latest, err := db.GetLatestVersion()
		if err != nil {
			logrus.WithError(err).Error("cannot get latest version")
			return
		}
		if err := db.LoadVersion(latest); err != nil {
			logrus.WithField("version", latest).WithError(err).Error("cannot load version")
			return
		}

		for key, order := range orderData {
			if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				b.Fatalf("cannot set key: %v", err)
			}
		}
		b.StartTimer()
		_, _, err = db.Commit()
		if err != nil {
			b.Fatalf("cannot commit: %v", err)
		}
		db.Close()
		b.StopTimer()
		size, _ := testsuite.GetDirSize(dir)
		b.Logf("dir size: %s", size.String())
		b.StartTimer()
	}
}
