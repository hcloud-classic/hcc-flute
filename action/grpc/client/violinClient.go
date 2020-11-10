package client

import (
	"context"
	"google.golang.org/grpc"
	"hcc/flute/action/grpc/pb/rpcviolin"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"strconv"
	"time"
)

var violinConn *grpc.ClientConn

func initViolin() error {
	var err error

	addr := config.Violin.ServerAddress + ":" + strconv.FormatInt(config.Violin.ServerPort, 10)
	violinConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.violin = rpcviolin.NewViolinClient(violinConn)
	logger.Logger.Println("gRPC violin client ready")

	return nil
}

func closeViolin() {
	_ = violinConn.Close()
}

// GetServerList : Get list of the server
func (rc *RPCClient) GetServerList(in *rpcviolin.ReqGetServerList) (*rpcviolin.ResGetServerList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetServerList, err := rc.violin.GetServerList(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetServerList, nil
}

// UpdateServer : Update infos of the server
func (rc *RPCClient) UpdateServer(in *rpcviolin.ReqUpdateServer) (*rpcviolin.ResUpdateServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resUpdateServer, err := rc.violin.UpdateServer(ctx, in)
	if err != nil {
		return nil, err
	}

	return resUpdateServer, nil
}
