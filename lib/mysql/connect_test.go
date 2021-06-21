package mysql

import (
	"hcc/flute/lib/config"
	"hcc/flute/lib/errors"
	"hcc/flute/lib/logger"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := logger.Init()
	if err != nil {
		errors.SetErrLogger(logger.Logger)
		errors.NewHccError(errors.FluteInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(logger.Logger)
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
