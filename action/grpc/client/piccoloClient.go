package client

import (
	"context"
	"google.golang.org/grpc"
	"hcc/flute/action/grpc/errconv"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"time"
)

var piccoloConn *grpc.ClientConn

func initPiccolo() error {
	var err error

	addr := config.Piccolo.ServerAddress + ":" + strconv.FormatInt(config.Piccolo.ServerPort, 10)
	logger.Logger.Println("Trying to connect to piccolo module (" + addr + ")")

	for i := 0; i < int(config.Piccolo.ConnectionRetryCount); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Piccolo.ConnectionTimeOutMs)*time.Millisecond)
		piccoloConn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logger.Logger.Println("Failed to connect piccolo module (" + addr + "): " + err.Error())
			logger.Logger.Println("Re-trying to connect to piccolo module (" +
				strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Piccolo.ConnectionRetryCount)) + ")")

			cancel()
			continue
		}

		RC.piccolo = pb.NewPiccoloClient(piccoloConn)
		logger.Logger.Println("gRPC client connected to piccolo module")

		cancel()
		return nil
	}

	hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "initPiccolo(): retry count exceeded to connect piccolo module")).Stack()
	return (*hccErrStack)[0].ToError()
}

func closePiccolo() {
	_ = piccoloConn.Close()
}

// WriteServerAlarm : Write server alarms to the database
func (rc *RPCClient) WriteServerAlarm(serverUUID string, reason string, detail string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Piccolo.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	_, err := rc.piccolo.WriteServerAlarm(ctx, &pb.ReqWriteServerAlarm{
		ServerUUID: serverUUID,
		ServerAlarm: &pb.ServerAlarm{
			Reason: reason,
			Detail: detail,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetGroupList : Get list of the group
func (rc *RPCClient) GetGroupList(_ *pb.Empty) (*pb.ResGetGroupList, *hcc_errors.HccErrorStack) {
	var errStack *hcc_errors.HccErrorStack

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Piccolo.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetGroupList, err := rc.piccolo.GetGroupList(ctx, &pb.Empty{})
	if err != nil {
		hccErrStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.FluteGrpcRequestError, "GetGroupList(): "+err.Error()))
		return nil, hccErrStack
	}
	if es := resGetGroupList.GetHccErrorStack(); es != nil {
		errStack = errconv.GrpcStackToHcc(es)
	}

	return resGetGroupList, errStack
}
