package daoext

import (
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
)

// ReadNodeNum : Get count of nodes from database.
func ReadNodeNum(in *pb.ReqGetNodeNum) (*pb.ResGetNodeNum, uint64, string) {
	var resNodeNum pb.ResGetNodeNum
	var nodeNr int64
	var groupID = in.GetGroupID()

	sql := "select count(*) from node where available = 1"
	if groupID != 0 {
		sql = "select count(*) from node where available = 1 and group_id = " + strconv.Itoa(int(groupID))
	}
	row := mysql.Db.QueryRow(sql)
	err := mysql.QueryRowScan(row, &nodeNr)
	if err != nil {
		errStr := "ReadNodeNum(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.FluteSQLNoResult, errStr
		}
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	resNodeNum.Num = nodeNr

	return &resNodeNum, 0, ""
}
