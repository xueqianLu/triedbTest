package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/triedbtest/cosmos"
	"github.com/xueqianLu/triedbtest/testsuite"
	"time"
)

func testCosmos(count int, dir string) error {
	rawdb, err := cosmos.NewRawDB(dir, true)
	if err != nil {
		logrus.WithError(err).Error("cannot create raw db")
		return err
	}
	defer rawdb.Close()

	db := cosmos.NewIAVL(rawdb)
	_, orderData := testsuite.GenerateAccount(count)
	latest, err := db.GetLatestVersion()
	if err != nil {
		logrus.WithError(err).Error("cannot get latest version")
		return err
	}
	if err := db.LoadVersion(latest); err != nil {
		logrus.WithField("version", latest).WithError(err).Error("cannot load version")
		return err
	}
	t1 := time.Now()

	for key, order := range orderData {
		if err := db.Set([]byte(fmt.Sprintf("ux-%s", key)), order); err != nil {
			logrus.WithError(err).Error("cannot set key")
			return err
		}
	}
	t2 := time.Now()
	logrus.WithFields(logrus.Fields{
		"stage": "write data to tree",
		"cost":  t2.Sub(t1).String(),
	})

	_, _, err = db.Commit()
	if err != nil {
		logrus.WithError(err).Error("cannot commit")
		return err
	}
	db.Close()
	t3 := time.Now()
	logrus.WithFields(logrus.Fields{
		"stage": "tree commit",
		"cost":  t3.Sub(t2).String(),
	})
	size, _ := testsuite.GetDirSize(dir)
	logrus.WithFields(logrus.Fields{
		"stage":    "total",
		"cost":     t3.Sub(t1).String(),
		"dir size": size.String(),
	})
	return nil
}
