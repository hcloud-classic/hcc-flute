package grpcsrv

import (
	"context"
	pb "hcc/flute/action/grpc/rpcflute"
	"hcc/flute/dao"
)

type fluteServer struct {
	pb.UnimplementedFluteServer
}

func (s *fluteServer) OnOffNode(_ context.Context, in *pb.ReqNodePowerControl) (*pb.ResNodePowerControl, error) {
	result, err := dao.NodePowerControl(in)
	if err != nil {
		return nil, err
	}

	return &pb.ResNodePowerControl{Result: result}, nil
}
