package uuidgen

import (
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/syscheck"
	"testing"
)

func Test_UUIDgen(t *testing.T) {
	if !syscheck.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := mysql.Prepare()
	if err != nil {
		return
	}
	defer func() {
		_ = mysql.Db.Close()
	}()

	_, err = UUIDgen()
	if err != nil {
		t.Fatal("Failed to generate uuid!")
	}
}
