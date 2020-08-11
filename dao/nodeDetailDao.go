package dao

import (
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
)

func ReadNodeDetail(args map[string]interface{}) (interface{}, error) {
	var nodeDetail model.NodeDetail
	var err error

	nodeUUID := args["node_uuid"].(string)
	var cpuModel string
	var cpuProcessors int
	var cpuThreads int

	sql := "select * from node_detail where node_uuid = ?"
	err = mysql.Db.QueryRow(sql, nodeUUID).Scan(
		&nodeUUID,
		&cpuModel,
		&cpuProcessors,
		&cpuThreads)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	nodeDetail.NodeUUID = nodeUUID
	nodeDetail.CPUModel = cpuModel
	nodeDetail.CPUProcessors = cpuProcessors
	nodeDetail.CPUThreads = cpuThreads

	return nodeDetail, nil
}
