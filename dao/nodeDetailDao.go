package dao

import (
	"hcc/flute/daoext"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

// CreateNodeDetail : Create detail infos of the node
func CreateNodeDetail(in *pb.ReqCreateNodeDetail) (*pb.NodeDetail, uint64, string) {
	reqNodeDetail := in.GetNodeDetail()
	if reqNodeDetail == nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNodeDetail(): nodeDetail is nil"
	}

	nodeUUID := reqNodeDetail.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	nodeDetailData := reqNodeDetail.GetNodeDetailData()
	nodeDetailDataOk := len(nodeDetailData) != 0
	nicDetailData := reqNodeDetail.GetNicDetailData()
	nicDetailDataOk := len(nicDetailData) != 0

	if !nodeUUIDOk || !nodeDetailDataOk || !nicDetailDataOk {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNodeDetail(): need nodeUUID and nodeDetailData, nicDetailData arguments"
	}

	nodeDetail := pb.NodeDetail{
		NodeUUID:       nodeUUID,
		NodeDetailData: reqNodeDetail.NodeDetailData,
		NicDetailData:  reqNodeDetail.NicDetailData,
	}

	sql := "insert into node_detail(node_uuid, node_detail_data, nic_detail_data) values (?, ?, ?)"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "CreateNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(nodeDetail.NodeUUID, nodeDetail.NodeDetailData)
	if err != nil {
		errStr := "CreateNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	return &nodeDetail, 0, ""
}

// UpdateNodeDetail : Update detail infos of the node
func UpdateNodeDetail(in *pb.ReqUpdateNodeDetail) (*pb.NodeDetail, uint64, string) {
	reqNodeDetail := in.GetNodeDetail()
	if reqNodeDetail == nil {
		return nil, hcc_errors.FluteGrpcRequestError, "UpdateNodeDetail(): nodeDetail is nil"
	}

	nodeUUID := reqNodeDetail.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	nodeDetailData := reqNodeDetail.GetNodeDetailData()
	nodeDetailDataOk := len(nodeDetailData) != 0
	nicDetailData := reqNodeDetail.GetNicDetailData()
	nicDetailDataOk := len(nicDetailData) != 0

	if !nodeUUIDOk {
		return nil, hcc_errors.FluteGrpcRequestError, "UpdateNodeDetail(): need a nodeUUID argument"
	}

	sql := "update node_detail set"
	var updateSet = ""

	if nodeDetailDataOk {
		updateSet += " node_detail_data = '" + nodeDetailData + "', "
	}
	if nicDetailDataOk {
		updateSet += " nic_detail_data = '" + nicDetailData + "', "
	}

	sql += updateSet[0:len(updateSet)-2] + " where node_uuid = ?"

	logger.Logger.Println("update_node_detail sql : ", sql)

	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "UpdateNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err2 := stmt.Exec(nodeUUID)
	if err2 != nil {
		errStr := "UpdateNodeDetail(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	defer func() {
		_ = stmt.Close()
	}()

	nodeDetail, errCode, errStr := daoext.ReadNodeDetail(nodeUUID)
	if errCode != 0 {
		logger.Logger.Println("UpdateNode(): " + errStr)
	}

	return nodeDetail, 0, ""
}

// DeleteNodeDetail : Delete detail infos of the node
func DeleteNodeDetail(in *pb.ReqDeleteNodeDetail) (*pb.NodeDetail, uint64, string) {
	nodeUUID := in.GetNodeUUID()
	nodeUUIDOk := len(nodeUUID) != 0
	if !nodeUUIDOk {
		return nil, hcc_errors.FluteGrpcRequestError, "DeleteNodeDetail(): need a nodeUUID argument"
	}

	nodeDetail, errCode, errText := daoext.ReadNodeDetail(nodeUUID)
	if errCode != 0 {
		return nil, hcc_errors.FluteGrpcRequestError, "DeleteNodeDetail(): " + errText
	}

	sql := "delete from node_detail where node_uuid = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "DeleteNodeDetail(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err2 := stmt.Exec(nodeUUID)
	if err2 != nil {
		errStr := "DeleteNodeDetail(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	logger.Logger.Println(result.RowsAffected())

	return nodeDetail, 0, ""
}
