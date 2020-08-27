package mysql

import (
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/syscheck"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := syscheck.CheckRoot()
	if err != nil {
		t.Fatal("Failed to get root permission!")
	}

	err = logger.Init()
	if err != nil {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Init()

	err = Init()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = Db.Close()
	}()
}
