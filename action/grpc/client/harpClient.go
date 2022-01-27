package client

import (
	"context"
	"google.golang.org/grpc"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"time"
)

var harpConn *grpc.ClientConn

func initHarp() error {
	var err error

	addr := config.Harp.ServerAddress + ":" + strconv.FormatInt(config.Harp.ServerPort, 10)
	logger.Logger.Println("Trying to connect to harp module (" + addr + ")")

	for i := 0; i < int(config.Harp.ConnectionRetryCount); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Harp.ConnectionTimeOutMs)*time.Millisecond)
		harpConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect harp module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to harp module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Harp.ConnectionRetryCount)) + ")")

			cancel()
			continue
		}

		RC.harp = pb.NewHarpClient(harpConn)
		logger.Logger.Println("gRPC client connected to harp module")

		cancel()
		return nil
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "initHarp(): retry count exceeded to connect harp module")).Stack()
	return (*hccErrStack)[0].ToError()
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
