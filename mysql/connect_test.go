package mysql

import (
	"hcloud-flute/checkroot"
	"hcloud-flute/config"
	"hcloud-flute/logger"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	if !checkroot.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer logger.FpLog.Close()

	config.Parser()

	err := Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer Db.Close()
}
