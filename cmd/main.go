package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/triedbtest/cosmos"
	"github.com/xueqianLu/triedbtest/testsuite"
	"math/big"
	"path/filepath"
)

func main() {
	dir := filepath.Join("./", "data")
	verifyUser := &types.StateAccount{
		CodeHash: []byte{0x1, 0x2, 0x3},
		Root:     common.HexToHash("0x123"),
		Nonce:    10018,
		Balance:  big.NewInt(111110),
	}
	verifyUserAddr := "122222222111111122"

	snapSet := testsuite.NewSnapshotSet()
	rawdb, err := cosmos.NewRawDB(dir, true)
	if err != nil {
		logrus.WithError(err).Error("cannot create raw db")
		return
	}
	defer rawdb.Close()

	for i := 0; i < 100; i++ {
		tree := cosmos.NewIAVL(rawdb)
		latest, err := tree.GetLatestVersion()
		if err != nil {
			logrus.WithError(err).Error("cannot get latest version")
			return
		}
		if err := tree.LoadVersion(latest); err != nil {
			logrus.WithField("version", latest).WithError(err).Error("cannot load version")
			return
		}
		_, orderData := testsuite.GenerateAccount(200000)
		verifyUser.Balance = new(big.Int).Add(verifyUser.Balance, big.NewInt(int64(i)))
		orderData[verifyUserAddr] = testsuite.AccountData(verifyUser)

		for key, order := range orderData {
			if err := tree.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
				logrus.WithError(err).Error("cannot set iavl trie")
				return
			}
		}
		_, newversion, err := tree.Commit()
		if err != nil {
			logrus.WithError(err).Error("cannot commit iavl trie")
			return
		}

		if err := tree.Close(); err != nil {
			logrus.WithError(err).Error("cannot close iavl trie")
		}
		snapSet.AddSnapshot(common.BytesToHash(big.NewInt(int64(newversion)).Bytes()), *verifyUser)
	}
	_, failed := snapSet.RangeSnapshot(func(sp testsuite.Snapshot) bool {
		vdb := cosmos.NewIAVL(rawdb)
		version := big.NewInt(0).SetBytes(sp.Root.Bytes()).Int64()
		d, err := vdb.Get(uint64(version), []byte(fmt.Sprintf("ux-%s", verifyUserAddr)))
		if err != nil {
			logrus.WithError(err).Error("cannot get order")
		}
		var order types.StateAccount
		if err := rlp.DecodeBytes(d, &order); err != nil {
			logrus.WithError(err).Error("cannot decode order")
		}
		if order.Balance.Cmp(sp.Account.Balance) != 0 {
			logrus.Errorf("balance mismatch: have %v, want %v", order.Balance, sp.Account.Balance)
		}
		vdb.Close()
		return true
	})
	if failed > 0 {
		logrus.Errorf("failed to verify %d snapshots", failed)
	} else {
		logrus.Info("all snapshots verified successfully")
	}
}
