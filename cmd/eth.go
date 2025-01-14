package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/triedbtest/ethtrie"
	"github.com/xueqianLu/triedbtest/testsuite"
	"time"
)

func testEth(count int, dir string) error {
	db := ethtrie.GetTrieDb(dir, true)
	defer db.Close()

	_, orderData := testsuite.GenerateAccount(count)
	tdb := trie.NewDatabase(db)
	// open tree, and set commit data to it.
	tree, err := trie.New(common.Hash{}, common.Hash{}, tdb)
	if err != nil {
		logrus.WithError(err).Error("cannot create trie")
		return err
	}
	t1 := time.Now()
	for key, order := range orderData {
		if err := tree.TryUpdate([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
			logrus.WithError(err).Error("cannot update trie")
			return err
		}
	}
	t2 := time.Now()
	logrus.WithFields(logrus.Fields{
		"stage": "write data to tree",
		"cost":  t2.Sub(t1).String(),
	}).Info("time info")
	merged := trie.NewMergedNodeSet()
	newroot, nodes, err := tree.Commit(true)
	if err != nil {
		logrus.WithError(err).Error("cannot commit trie")
		return err
	}
	if err = merged.Merge(nodes); err != nil {
		logrus.WithError(err).Error("cannot merge node")
		return err
	}

	if err = tdb.Update(merged); err != nil {
		logrus.WithError(err).Error("cannot update trie")
		return err
	}
	if err = tdb.Commit(newroot, false, nil); err != nil {
		logrus.WithError(err).Error("cannot commit trie")
		return err
	}
	t3 := time.Now()

	logrus.WithFields(logrus.Fields{
		"stage": "tree commit",
		"cost":  t3.Sub(t2).String(),
	}).Info("time info")
	size, _ := testsuite.GetDirSize(dir)
	logrus.WithFields(logrus.Fields{
		"stage":    "total",
		"cost":     t3.Sub(t1).String(),
		"dir size": size.String(),
	}).Info("time info")
	return nil
}
