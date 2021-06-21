package daoext

import (
	pb "hcc/flute/action/grpc/pb/rpcflute"
	hccerr "hcc/flute/lib/errors"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"strings"
)

// ReadNodeNum : Get count of nodes from database.
func ReadNodeNum() (*pb.ResGetNodeNum, uint64, string) {
	var resNodeNum pb.ResGetNodeNum
	var nodeNr int64

	sql := "select count(*) from node where available = 1"
	err := mysql.Db.QueryRow(sql).Scan(&nodeNr)
	if err != nil {
		errStr := "ReadNodeNum(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hccerr.FluteSQLNoResult, errStr
		}
		return nil, hccerr.FluteSQLOperationFail, errStr
	}
	resNodeNum.Num = nodeNr

	return &resNodeNum, 0, ""
}
