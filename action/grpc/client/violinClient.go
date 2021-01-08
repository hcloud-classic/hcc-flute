package client

import (
	"context"
	"github.com/hcloud-classic/hcc_errors"
	"github.com/hcloud-classic/pb"
	"google.golang.org/grpc"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"strconv"
	"time"
)

var violinConn *grpc.ClientConn

func initViolin() error {
	var err error

	addr := config.Violin.ServerAddress + ":" + strconv.FormatInt(config.Violin.ServerPort, 10)
	logger.Logger.Println("Trying to connect to violin module (" + addr + ")")

	for i := 0; i < int(config.Violin.ConnectionRetryCount); i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Violin.ConnectionTimeOutMs)*time.Millisecond)
		violinConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect violin module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to violin module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Violin.ConnectionRetryCount)) + ")")
			continue
		}

		RC.violin = pb.NewViolinClient(violinConn)
		logger.Logger.Println("gRPC client connected to violin module")

		return nil
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "initViolin(): retry count exceeded to connect violin module")).Stack()
	return (*hccErrStack)[0].ToError()
}

func closeViolin() {
	_ = violinConn.Close()
}

// GetServerList : Get list of the server
func (rc *RPCClient) GetServerList(in *pb.ReqGetServerList) (*pb.ResGetServerList, error) {
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
func (rc *RPCClient) UpdateServer(in *pb.ReqUpdateServer) (*pb.ResUpdateServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resUpdateServer, err := rc.violin.UpdateServer(ctx, in)
	if err != nil {
		return nil, err
	}

	return resUpdateServer, nil
}
