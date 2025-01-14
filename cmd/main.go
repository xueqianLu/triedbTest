package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var (
	caseFlag  = flag.String("case", "eth", "db type to test")
	datacount = flag.Int("count", 200000, "data count in one test")
	testTimes = flag.Int("N", 100, "test times")
)

func main() {
	flag.Parse()

	dir := filepath.Join("./", "data-"+*caseFlag)
	for i := 0; i < *testTimes; i++ {
		var err error
		switch *caseFlag {
		case "eth":
			err = testEth(*datacount, dir)
		default:
			err = testCosmos(*datacount, dir)
		}
		if err != nil {
			logrus.WithField("test idx", i).WithError(err).Error("test failed")
		}
	}
	logrus.Info("test finished")
}
