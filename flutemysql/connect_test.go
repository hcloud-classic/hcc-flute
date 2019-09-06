package cellomysql

import (
	"GraphQL_Cello/cellocheckroot"
	"GraphQL_Cello/cellologger"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	if !cellocheckroot.CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}

	if !cellologger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer cellologger.FpLog.Close()

	err := Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer Db.Close()
}
