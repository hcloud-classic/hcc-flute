package dao

import (
	pb "hcc/flute/action/grpc/pb/rpcflute"
	hccerr "hcc/flute/lib/errors"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"strings"
)

// ReadNodeDetail : Get detail infos of the node
func ReadNodeDetail(nodeUUID string) (*pb.NodeDetail, uint64, string) {
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
		errStr := "ReadNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hccerr.FluteSQLNoResult, errStr
		}
		return nil, hccerr.FluteSQLOperationFail, errStr
	}

	nodeDetail.NodeUUID = nodeUUID
	nodeDetail.CPUModel = cpuModel
	nodeDetail.CPUProcessors = int32(cpuProcessors)
	nodeDetail.CPUThreads = int32(cpuThreads)

	return &nodeDetail, 0, ""
}

// CreateNodeDetail : Create detail infos of the node
func CreateNodeDetail(in *pb.ReqCreateNodeDetail) (*pb.NodeDetail, uint64, string) {
	reqNodeDetail := in.GetNodeDetail()
	if reqNodeDetail == nil {
		return nil, hccerr.FluteGrpcRequestError, "CreateNodeDetail(): nodeDetail is nil"
	}

	nodeUUID := reqNodeDetail.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	if !nodeUUIDOk {
		return nil, hccerr.FluteGrpcRequestError, "CreateNodeDetail(): need a nodeUUID argument"
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
		errStr := "CreateNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(nodeDetail.NodeUUID, nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads)
	if err != nil {
		errStr := "CreateNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}

	return &nodeDetail, 0, ""
}

// DeleteNodeDetail : Delete detail infos of the node
func DeleteNodeDetail(in *pb.ReqDeleteNodeDetail) (string, uint64, string) {
	nodeUUID := in.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	if !nodeUUIDOk {
		return "", hccerr.FluteGrpcRequestError, "DeleteNodeDetail(): need a nodeUUID argument"
	}

	sql := "delete from node_detail where node_uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "DeleteNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.FluteSQLOperationFail, errStr
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err2 := stmt.Exec(nodeUUID)
	if err2 != nil {
		errStr := "DeleteNodeDetail(): " + err2.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.FluteSQLOperationFail, errStr
	}
	logger.Logger.Println(result.RowsAffected())

	return nodeUUID, 0, ""
}
