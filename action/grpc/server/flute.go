package server

import (
	"context"
	"hcc/flute/action/grpc/errconv"
	"hcc/flute/dao"
	"hcc/flute/daoext"
	"hcc/flute/lib/logger"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

type fluteServer struct {
	pb.UnimplementedFluteServer
}

func returnNode(node *pb.Node) *pb.Node {
	return &pb.Node{
		UUID:            node.UUID,
		GroupID:         node.GroupID,
		NodeNum:         node.NodeNum,
		NodeIP:          node.NodeIP,
		ServerUUID:      node.ServerUUID,
		BmcMacAddr:      node.BmcMacAddr,
		BmcIP:           node.BmcIP,
		BmcIPSubnetMask: node.BmcIPSubnetMask,
		PXEMacAddr:      node.PXEMacAddr,
		Status:          node.Status,
		CPUCores:        node.CPUCores,
		Memory:          node.Memory,
		NicSpeedMbps:    node.NicSpeedMbps,
		Description:     node.Description,
		RackNumber:      node.RackNumber,
		ChargeCPU:       node.ChargeCPU,
		ChargeMemory:    node.ChargeMemory,
		ChargeNIC:       node.ChargeNIC,
		Active:          node.Active,
		CreatedAt:       node.CreatedAt,
	}
}

func returnNodeDetail(nodeDetail *pb.NodeDetail) *pb.NodeDetail {
	return &pb.NodeDetail{
		NodeUUID:       nodeDetail.NodeUUID,
		NodeDetailData: nodeDetail.NodeDetailData,
	}
}

func (s *fluteServer) CreateNode(_ context.Context, in *pb.ReqCreateNode) (*pb.ResCreateNode, error) {
	logger.Logger.Println("Request received: CreateNode()")

	node, errCode, errStr := dao.CreateNode(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResCreateNode{Node: &pb.Node{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreateNode{Node: returnNode(node)}, nil
}

func (s *fluteServer) GetNode(_ context.Context, in *pb.ReqGetNode) (*pb.ResGetNode, error) {
	logger.Logger.Println("Request received: GetNode()")

	node, errCode, errStr := dao.ReadNode(in.GetUUID())
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetNode{Node: &pb.Node{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetNode{Node: returnNode(node)}, nil
}

func (s *fluteServer) GetNodeList(_ context.Context, in *pb.ReqGetNodeList) (*pb.ResGetNodeList, error) {
	logger.Logger.Println("Request received: GetNodeList()")

	nodeList, errCode, errStr := dao.ReadNodeList(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetNodeList{Node: []*pb.Node{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return nodeList, nil
}

func (s *fluteServer) GetNodeNum(_ context.Context, in *pb.ReqGetNodeNum) (*pb.ResGetNodeNum, error) {
	logger.Logger.Println("Request received: GetNodeNum()")

	nodeNum, errCode, errStr := daoext.ReadNodeNum(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetNodeNum{Num: 0, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil

	}

	return nodeNum, nil
}

func (s *fluteServer) UpdateNode(_ context.Context, in *pb.ReqUpdateNode) (*pb.ResUpdateNode, error) {
	logger.Logger.Println("Request received: UpdateNode()")

	updateNode, errCode, errStr := dao.UpdateNode(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResUpdateNode{Node: &pb.Node{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResUpdateNode{Node: updateNode}, nil
}

func (s *fluteServer) DeleteNode(_ context.Context, in *pb.ReqDeleteNode) (*pb.ResDeleteNode, error) {
	logger.Logger.Println("Request received: DeleteNode()")

	deleteNode, errCode, errStr := dao.DeleteNode(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResDeleteNode{Node: &pb.Node{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResDeleteNode{Node: deleteNode}, nil
}

func (s *fluteServer) NodePowerControl(_ context.Context, in *pb.ReqNodePowerControl) (*pb.ResNodePowerControl, error) {
	logger.Logger.Println("Request received: NodePowerControl()")

	result, errCode, errStr := dao.NodePowerControl(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResNodePowerControl{Result: []string{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResNodePowerControl{Result: result}, nil
}

func (s *fluteServer) GetNodePowerState(_ context.Context, in *pb.ReqNodePowerState) (*pb.ResNodePowerState, error) {
	logger.Logger.Println("Request received: GetNodePowerState()")

	result, errCode, errStr := dao.GetNodePowerState(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResNodePowerState{Result: "", HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResNodePowerState{Result: result}, nil
}

func (s *fluteServer) CreateNodeDetail(_ context.Context, in *pb.ReqCreateNodeDetail) (*pb.ResCreateNodeDetail, error) {
	logger.Logger.Println("Request received: CreateNodeDetail()")

	nodeDetail, errCode, errStr := dao.CreateNodeDetail(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResCreateNodeDetail{NodeDetail: &pb.NodeDetail{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResCreateNodeDetail{NodeDetail: returnNodeDetail(nodeDetail)}, nil
}

func (s *fluteServer) GetNodeDetail(_ context.Context, in *pb.ReqGetNodeDetail) (*pb.ResGetNodeDetail, error) {
	logger.Logger.Println("Request received: GetNodeDetail()")

	nodeDetail, errCode, errStr := daoext.ReadNodeDetail(in.GetNodeUUID())
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResGetNodeDetail{NodeDetail: &pb.NodeDetail{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetNodeDetail{NodeDetail: returnNodeDetail(nodeDetail)}, nil
}

func (s *fluteServer) UpdateNodeDetail(_ context.Context, in *pb.ReqUpdateNodeDetail) (*pb.ResUpdateNodeDetail, error) {
	logger.Logger.Println("Request received: UpdateNodeDetail()")

	nodeDetail, errCode, errStr := dao.UpdateNodeDetail(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResUpdateNodeDetail{NodeDetail: &pb.NodeDetail{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResUpdateNodeDetail{NodeDetail: returnNodeDetail(nodeDetail)}, nil
}

func (s *fluteServer) DeleteNodeDetail(_ context.Context, in *pb.ReqDeleteNodeDetail) (*pb.ResDeleteNodeDetail, error) {
	logger.Logger.Println("Request received: DeleteNodeDetail()")

	nodeDetail, errCode, errStr := dao.DeleteNodeDetail(in)
	if errCode != 0 {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(errCode, errStr))
		return &pb.ResDeleteNodeDetail{NodeDetail: &pb.NodeDetail{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResDeleteNodeDetail{NodeDetail: nodeDetail}, nil
}
