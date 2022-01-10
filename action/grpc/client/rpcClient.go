package client

import (
	"innogrid.com/hcloud-classic/pb"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	horn    pb.HornClient
	violin  pb.ViolinClient
	piccolo pb.PiccoloClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// Init : Initialize clients of gRPC
func Init() error {
	err := initHorn()
	if err != nil {
		return err
	}

	err = initViolin()
	if err != nil {
		return err
	}

	err = initPiccolo()
	if err != nil {
		return err
	}

	return nil
}

// End : Close connections of gRPC clients
func End() {
	closePiccolo()
	closeViolin()
	closeHorn()
}
