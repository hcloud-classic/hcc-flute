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

func CreateNodeDetail(args map[string]interface{}) (interface{}, error) {
	var err error

	nodeUUID, nodeUUIDOk := args["node_uuid"].(string)
	if !nodeUUIDOk {
		return nil, err
	}

	nodeDetail := model.NodeDetail{
		NodeUUID:      nodeUUID,
		CPUModel:      args["cpu_model"].(string),
		CPUProcessors: args["cpu_processors"].(int),
		CPUThreads:    args["cpu_threads"].(int),
	}

	sql := "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(nodeDetail.NodeUUID, nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return nodeDetail, nil
}

func DeleteNodeDetail(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, ok := args["node_uuid"].(string)
	if ok {
		sql := "delete from node_detail where node_uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}

		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(requestedUUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.RowsAffected())

		return requestedUUID, nil
	}

	return requestedUUID, err
}
