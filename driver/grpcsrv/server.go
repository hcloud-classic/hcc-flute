package grpcsrv

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hcc/flute/action/grpc/rpcflute"
)

const (
	port = ":50051"
)

type server struct {
	rpcflute.UnimplementedFluteServer
}

func (s *server) OnOffNode(ctx context.Context, in *rpcflute.ReqOnOffNode) (*rpcflute.ResOnOffNode, error) {
	fmt.Println(in.Nodes)
	node := rpcflute.Node{Uuid: "uuid"}
	nodes := []*rpcflute.Node{&node}
	return &rpcflute.ResOnOffNode{Nodes: nodes}, nil
}

func InitGRPC() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	rpcflute.RegisterFluteServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func CleanGRPC() {

}
