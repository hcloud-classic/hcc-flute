package client

import (
	"hcc/flute/action/grpc/pb/rpcviolin"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	violin rpcviolin.ViolinClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// Init : Initialize clients of gRPC
func Init() error {
	err := initViolin()
	if err != nil {
		return err
	}

	return nil
}

// End : Close connections of gRPC clients
func End() {
	closeViolin()
}
