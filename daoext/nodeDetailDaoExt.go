package daoext

import (
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strings"
)

// ReadNodeDetail : Get detail infos of the node
func ReadNodeDetail(nodeUUID string) (*pb.NodeDetail, uint64, string) {
	var nodeDetail pb.NodeDetail

	var nodeDetailData string
	var nicDetailData string

	sql := "select * from node_detail where node_uuid = ?"
	row := mysql.Db.QueryRow(sql, nodeUUID)
	err := mysql.QueryRowScan(row,
		&nodeUUID,
		&nodeDetailData,
		&nicDetailData)
	if err != nil {
		errStr := "ReadNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.FluteSQLNoResult, errStr
		}
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	nodeDetail.NodeUUID = nodeUUID
	nodeDetail.NodeDetailData = nodeDetailData
	nodeDetail.NicDetailData = nicDetailData

	return &nodeDetail, 0, ""
}
