package grpccli

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"hcc/flute/action/grpc/rpcharp"
)

const (
	address = "localhost:50052"
)

var harp rpcharp.HarpClient
var conn grpc.ClientConn

func InitGRPCHarp() error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Connect Harp failed: %v", err)
		return err
	}

	harp = rpcharp.NewHarpClient(conn)

	return nil
}

func CleanGRPCHarp() {

	conn.Close()
}

func GetSubnet(req *rpcharp.ReqGetSubnet) (*rpcharp.ResGetSubnet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second /* 10 secs */)
	defer cancel()
	r, err := harp.GetSubnet(ctx, req)
	if err != nil {
		return r, err
	}
	return r, nil
}
