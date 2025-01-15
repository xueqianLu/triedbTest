package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/triedbtest/cosmos"
	"github.com/xueqianLu/triedbtest/testsuite"
	"time"
)

func testCosmos(orderData map[string][]byte, idx int, count int, dir string) error {
	loger := logrus.WithField("idx", idx)
	rawdb, err := cosmos.NewRawDB(dir, true)
	if err != nil {
		loger.WithError(err).Error("cannot create raw db")
		return err
	}
	defer rawdb.Close()

	db := cosmos.NewIAVL(rawdb)
	latest, err := db.GetLatestVersion()
	if err != nil {
		loger.WithError(err).Error("cannot get latest version")
		return err
	}
	if err := db.LoadVersion(latest); err != nil {
		loger.WithField("version", latest).WithError(err).Error("cannot load version")
		return err
	}
	t1 := time.Now()

	for key, order := range orderData {
		if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
			loger.WithError(err).Error("cannot set key")
			return err
		}
	}
	t2 := time.Now()
	loger.WithFields(logrus.Fields{
		"stage": "write data to tree",
		"cost":  t2.Sub(t1).String(),
	}).Info("time info")

	_, _, err = db.Commit()
	if err != nil {
		loger.WithError(err).Error("cannot commit")
		return err
	}
	db.Close()
	t3 := time.Now()
	loger.WithFields(logrus.Fields{
		"stage": "tree commit",
		"cost":  t3.Sub(t2).String(),
	}).Info("time info")
	size, _ := testsuite.GetDirSize(dir)
	loger.WithFields(logrus.Fields{
		"stage":    "total",
		"cost":     t3.Sub(t1).String(),
		"dir size": size.String(),
	}).Info("time info")
	return nil
}
