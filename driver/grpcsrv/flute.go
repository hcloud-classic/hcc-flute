package grpcsrv

import (
	"context"
	pb "hcc/flute/action/grpc/rpcflute"
	"hcc/flute/dao"
)

type fluteServer struct {
	pb.UnimplementedFluteServer
}

func returnNode(node *pb.Node) *pb.Node {
	return &pb.Node{
		ServerUUID: node.ServerUUID,
		BmcMacAddr: node.BmcMacAddr,
		BmcIP: node.BmcIP,
		PXEMacAddr: node.PXEMacAddr,
		Status: node.Status,
		CPUCores: node.CPUCores,
		Memory: node.Memory,
		Description: node.Description,
		Active: node.Active,
		CreatedAt: node.CreatedAt,
	}
}

func returnNodeDetail(nodeDetail *pb.NodeDetail) *pb.NodeDetail {
	return &pb.NodeDetail{
		NodeUUID: nodeDetail.NodeUUID,
		CPUModel: nodeDetail.CPUModel,
		CPUProcessors: nodeDetail.CPUProcessors,
		CPUThreads: nodeDetail.CPUThreads,
	}
}

func (s *fluteServer) CreateNode(_ context.Context, in *pb.ReqCreateNode) (*pb.ResCreateNode, error) {
	node, err := dao.CreateNode(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateNode{Node: returnNode(node)}, nil
}

func (s *fluteServer) GetNode(_ context.Context, in *pb.ReqGetNode) (*pb.ResGetNode, error) {
	node, err := dao.ReadNode(in.GetUUID())
	if err != nil {
		return nil, err
	}

	return &pb.ResGetNode{Node: returnNode(node)}, nil
}

func (s *fluteServer) GetNodeList(_ context.Context, in *pb.ReqGetNodeList) (*pb.ResGetNodeList, error) {
	nodeList, err := dao.ReadNodeList(in)
	if err != nil {
		return nil, err
	}

	return nodeList, nil
}

func (s *fluteServer) GetNodeNum(_ context.Context, _ *pb.Empty) (*pb.ResGetNodeNum, error) {
	nodeNum, err := dao.ReadNodeNum()
	if err != nil {
		return nil, err
	}

	return nodeNum, nil
}

func (s *fluteServer) UpdateNode(_ context.Context, in *pb.ReqUpdateNode) (*pb.ResUpdateNode, error) {
	updateNode, err := dao.UpdateNode(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResUpdateNode{Node: updateNode}, nil
}

func (s *fluteServer) DeleteNode(_ context.Context, in *pb.ReqDeleteNode) (*pb.ResDeleteNode, error) {
	uuid, err := dao.DeleteNode(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResDeleteNode{UUID: uuid}, nil
}

func (s *fluteServer) NodePowerControl(_ context.Context, in *pb.ReqNodePowerControl) (*pb.ResNodePowerControl, error) {
	result, err := dao.NodePowerControl(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResNodePowerControl{Result: result}, nil
}

func (s *fluteServer) GetNodePowerState(_ context.Context, in *pb.ReqNodePowerState) (*pb.ResNodePowerState, error) {
	result, err := dao.GetNodePowerState(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResNodePowerState{Result: result}, nil
}

func (s *fluteServer) CreateNodeDetail(_ context.Context, in *pb.ReqCreateNodeDetail) (*pb.ResCreateNodeDetail, error) {
	nodeDetail, err := dao.CreateNodeDetail(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResCreateNodeDetail{NodeDetail: returnNodeDetail(nodeDetail)}, nil
}

func (s *fluteServer) GetNodeDetail(_ context.Context, in *pb.ReqGetNodeDetail) (*pb.ResGetNodeDetail, error) {
	nodeDetail, err := dao.ReadNodeDetail(in.GetNodeUUID())
	if err != nil {
		return nil, err
	}

	return &pb.ResGetNodeDetail{NodeDetail: returnNodeDetail(nodeDetail)}, nil
}

func (s *fluteServer) DeleteNodeDetail(_ context.Context, in *pb.ReqDeleteNodeDetail) (*pb.ResDeleteNodeDetail, error) {
	nodeUUID, err := dao.DeleteNodeDetail(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResDeleteNodeDetail{NodeUUID: nodeUUID}, nil
}
