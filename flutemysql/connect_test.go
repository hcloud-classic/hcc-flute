package flutemysql

import (
	"GraphQL_Flute/flutecheckroot"
	"GraphQL_Flute/flutelogger"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	if !flutecheckroot.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !flutelogger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer flutelogger.FpLog.Close()

	err := Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer Db.Close()
}
