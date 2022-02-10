package mysql

import (
	"hcc/flute/action/grpc/client"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"innogrid.com/hcloud-classic/hcc_errors"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {
	err := logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(logger.Logger)
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Init()

	err = client.Init()
	if err != nil {
		t.Fatal(err)
	}

	err = Init()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = Db.Close()
	}()
}
