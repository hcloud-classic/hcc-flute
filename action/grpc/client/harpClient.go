package client

import (
	"context"
	"google.golang.org/grpc"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"time"
)

var harpConn *grpc.ClientConn

func initHarp() error {
	var err error

	addr := config.Harp.ServerAddress + ":" + strconv.FormatInt(config.Harp.ServerPort, 10)
	harpConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.harp = pb.NewHarpClient(harpConn)
	logger.Logger.Println("gRPC harp client ready")

	return nil
}

func closeHarp() {
	_ = harpConn.Close()
}

// GetSubnetByServer : Get infos of the subnet
func (rc *RPCClient) GetSubnetByServer(serverUUID string) (*pb.ResGetSubnetByServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetSubnetByServer, err := rc.harp.GetSubnetByServer(ctx, &pb.ReqGetSubnetByServer{ServerUUID: serverUUID})
	if err != nil {
		return nil, err
	}

	return resGetSubnetByServer, nil
}
