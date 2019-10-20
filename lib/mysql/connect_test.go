package mysql

import (
<<<<<<< HEAD
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/syscheck"
=======
	"hcc/flute/lib/syscheck"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
>>>>>>> f41ff24 (Refactoring packages structure)
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
<<<<<<< HEAD
	err := syscheck.CheckRoot()
	if err != nil {
=======
	if !syscheck.CheckRoot() {
>>>>>>> f41ff24 (Refactoring packages structure)
		t.Fatal("Failed to get root permission!")
	}

	if !logger.Prepare() {
		t.Fatal("Failed to prepare logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

<<<<<<< HEAD
	err = Prepare()
=======
	err := Prepare()
>>>>>>> f41ff24 (Refactoring packages structure)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = Db.Close()
	}()
}
