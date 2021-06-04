package dao

import (
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strings"
)

// GetNodeUptime : Get the node's uptime info from the database
func GetNodeUptime(nodeUUID string, day string) (*pb.NodeUptime, uint64, string) {
	var nodeUptime pb.NodeUptime

	var groupID int64
	var uptimeMs int64

	sql := "select group_id, uptime_ms from node_uptime where node_uuid = ? and day = ?"
	row := mysql.Db.QueryRow(sql, nodeUUID, day)
	err := mysql.QueryRowScan(row,
		&groupID,
		&uptimeMs)
	if err != nil {
		errStr := "GetNodeUptime(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.FluteSQLNoResult, errStr
		}
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	nodeUptime.NodeUUID = nodeUUID
	nodeUptime.GroupID = groupID
	nodeUptime.Day = day
	nodeUptime.UptimeMs = uptimeMs

	return &nodeUptime, 0, ""
}
