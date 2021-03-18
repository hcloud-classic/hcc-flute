package errconv

import (
	errh "innogrid.com/hcloud-classic/hcc_errors"
	errg "innogrid.com/hcloud-classic/pb"
)

func GrpcToHcc(eg *errg.HccError) *errh.HccError {
	return errh.NewHccError(eg.GetErrCode(), eg.GetErrText())
}

func HccToGrpc(eh *errh.HccError) *errg.HccError {
	return &errg.HccError{ErrCode: eh.Code(), ErrText: eh.Text()}
}

func GrpcStackToHcc(esg *errg.HccErrorStack) *errh.HccErrorStack {
	errStack := errh.NewHccErrorStack()

	if errStack.Version() != esg.GetVersion() {
		errStack.IsMixed = true
	} else {
		errStack.IsMixed = false
	}

	for _, e := range esg.GetErrStack() {
		errStack.Push(errh.NewHccError(e.GetErrCode(), e.GetErrText()))
	}

	return errStack
}

func HccStackToGrpc(esh *errh.HccErrorStack) *errg.HccErrorStack {
	ges := new(errg.HccErrorStack)

	ges.Version = esh.Version()
	ges.IsMixed = esh.IsMixed

	for i := 0; i < esh.Len(); i++ {
		eh := esh.Pop()
		ge := HccToGrpc(eh)
		ges.ErrStack = append(ges.GetErrStack(), ge)
	}

	return ges
}
