package errconv

import (
	errh "github.com/hcloud-classic/hcc_errors"
	errg "github.com/hcloud-classic/pb"
)

// GrpcToHcc : Convert gRPC error code and text to HCCError
func GrpcToHcc(eg *errg.HccError) *errh.HccError {
	return errh.NewHccError(eg.GetErrCode(), eg.GetErrText())
}

// HccToGrpc : Convert HCCError to gRPC error
func HccToGrpc(eh *errh.HccError) *errg.HccError {
	return &errg.HccError{ErrCode: eh.Code(), ErrText: eh.Text()}
}

// GrpcStackToHcc : Convert gRPC error stack to HCCError stack
func GrpcStackToHcc(esg *[]*errg.HccError) *errh.HccErrorStack {
	errStack := errh.NewHccErrorStack()

	for i, e := range *esg {
		if i == 0 {
			continue
		}

		_ = errStack.Push(errh.NewHccError(e.GetErrCode(), e.GetErrText()))
	}

	return errStack
}

// HccStackToGrpc : Convert HCCError stack to gRPC error stack
func HccStackToGrpc(esh *errh.HccErrorStack) []*errg.HccError {
	ges := []*errg.HccError{}
	for i := 0; i <= esh.Len(); i++ {
		ge := &errg.HccError{ErrCode: (*esh.Stack())[i].Code(), ErrText: (*esh.Stack())[i].Text()}
		ges = append(ges, ge)
	}
	return ges
}
