package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/triedbtest/ethtrie"
	"github.com/xueqianLu/triedbtest/testsuite"
	"path/filepath"
)

var (
	caseFlag  = flag.String("case", "eth", "db type to test")
	dataType  = flag.String("data", "custom", "data type to test , custom or eth")
	dataSize  = flag.Int("size", 1, "custom data size, 1 is 32bytes")
	datacount = flag.Int("count", 200000, "data count in one test")
	testTimes = flag.Int("N", 100, "test times")
)

func main() {
	flag.Parse()

	dir := filepath.Join("./", "data-"+*caseFlag)

	ethdb := ethtrie.GetTrieDb(dir, true)
	defer ethdb.Close()

	orderdata := make(map[string][]byte)

	for i := 0; i < *testTimes; i++ {
		switch *dataType {
		case "custom":
			_, orderdata = testsuite.GenerateCustom(*datacount, *dataSize)
		default:
			_, orderdata = testsuite.GenerateAccount(*datacount)
		}
		var err error
		switch *caseFlag {
		case "eth":
			err = testEth(orderdata, ethdb, i, *datacount, dir)
		default:
			err = testCosmos(orderdata, i, *datacount, dir)
		}
		if err != nil {
			logrus.WithField("test idx", i).WithError(err).Error("test failed")
		}
	}
	logrus.Info("test finished")
}
