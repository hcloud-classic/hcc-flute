package errconv

import (
	errg "hcc/flute/action/grpc/pb/rpcmsgType"
	errh "hcc/flute/lib/errors"
)

func GrpcToHcc(eg *errg.HccError) *errh.HccError {
	return errh.NewHccError(eg.GetErrCode(), eg.GetErrText())
}

func HccToGrpc(eh *errh.HccError) *errg.HccError {
	return &errg.HccError{ErrCode: eh.ErrCode, ErrText: eh.ErrText}
}

func GrpcStackToHcc(esg *[]*errg.HccError) *errh.HccErrorStack {
	errStack := errh.NewHccErrorStack()

	for _, e := range *esg {
		errStack.Push(errh.NewHccError(e.GetErrCode(), e.GetErrText()))
	}

	hccErrStack := *errStack
	es := hccErrStack[1:]

	return &es
}

func HccStackToGrpc(esh *errh.HccErrorStack) []*errg.HccError {
	ges := []*errg.HccError{}
	for i := esh.Len(); i >= 0; i-- {
		ge := &errg.HccError{ErrCode: (*esh)[i].ErrCode, ErrText: (*esh)[i].ErrText}
		ges = append(ges, ge)
	}
	return ges
}
