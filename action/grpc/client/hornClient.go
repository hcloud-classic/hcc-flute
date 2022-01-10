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

var hornConn *grpc.ClientConn

func initHorn() error {
	var err error

	addr := config.Horn.ServerAddress + ":" + strconv.FormatInt(config.Horn.ServerPort, 10)
	logger.Logger.Println("Trying to connect to horn module (" + addr + ")")

	for i := 0; i < int(config.Horn.ConnectionRetryCount); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Horn.ConnectionTimeOutMs)*time.Millisecond)
		hornConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect horn module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to horn module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Horn.ConnectionRetryCount)) + ")")

			cancel()
			continue
		}

		RC.horn = pb.NewHornClient(hornConn)
		logger.Logger.Println("gRPC client connected to horn module")

		cancel()
		return nil
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.HarpInternalInitFail, "initHorn(): retry count exceeded to connect horn module")).Stack()
	return (*hccErrStack)[0].ToError()
}

func closeHorn() {
	_ = hornConn.Close()
}

// GetMYSQLDEncryptedPassword : Get encrypted password of mysqld
func (rc *RPCClient) GetMYSQLDEncryptedPassword() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Horn.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetMYSQLDEncryptedPassword, err := rc.horn.GetMYSQLDEncryptedPassword(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return resGetMYSQLDEncryptedPassword.EncryptedPassword, nil
}
