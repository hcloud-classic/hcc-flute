package dao

import (
	"errors"
	pb "hcc/flute/action/grpc/rpcflute"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
)

// ReadNodeDetail : Get detail infos of the node
func ReadNodeDetail(nodeUUID string) (*pb.NodeDetail, error) {
	var nodeDetail pb.NodeDetail

	var cpuModel string
	var cpuProcessors int
	var cpuThreads int

	sql := "select * from node_detail where node_uuid = ?"
	err := mysql.Db.QueryRow(sql, nodeUUID).Scan(
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
	nodeDetail.CPUProcessors = int32(cpuProcessors)
	nodeDetail.CPUThreads = int32(cpuThreads)

	return &nodeDetail, nil
}

// CreateNodeDetail : Create detail infos of the node
func CreateNodeDetail(in *pb.ReqCreateNodeDetail) (*pb.NodeDetail, error) {
	reqNodeDetail := in.GetNodeDetail()
	if reqNodeDetail == nil {
		return nil, errors.New("nodeDetail is nil")
	}

	nodeUUID := reqNodeDetail.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	if !nodeUUIDOk {
		return nil, errors.New("need a nodeUUID argument")
	}

	nodeDetail := pb.NodeDetail{
		NodeUUID:      nodeUUID,
		CPUModel:      reqNodeDetail.CPUModel,
		CPUProcessors: reqNodeDetail.CPUProcessors,
		CPUThreads:    reqNodeDetail.CPUThreads,
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

	return &nodeDetail, nil
}

// DeleteNodeDetail : Delete detail infos of the node
func DeleteNodeDetail(in *pb.ReqDeleteNodeDetail) (string, error) {
	nodeUUID := in.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	if !nodeUUIDOk {
		return "", errors.New("need a nodeUUID argument")
	}

	sql := "delete from node_detail where node_uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return "", err
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err2 := stmt.Exec(nodeUUID)
	if err2 != nil {
		logger.Logger.Println(err2)
		return "", err
	}
	logger.Logger.Println(result.RowsAffected())

	return nodeUUID, nil
}
