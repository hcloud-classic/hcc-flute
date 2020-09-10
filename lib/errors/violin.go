package errors

const (
	ViolinInternalInitFail                 = violin + internal + initFail
	ViolinInternalOperationFail            = violin + internal + operationFail
	ViolinInternalUUIDGenerationError      = violin + internal + UUIDGenerationError
	ViolinInternalTimeStampConversionError = violin + internal + timestampConversionError
	ViolinInternalCreateServerFailed       = violin + internal + createServerFailed
	ViolinInternalCreateServerRoutineError = violin + internal + createServerRoutineError
	ViolinInternalGetAvailableNodesError   = violin + internal + getAvailableNodesError
	ViolinInternalServerNodePresentError   = violin + internal + serverNodePresentError

	ViolinGrpcArgumentError = violin + grpc + argumentError
	ViolinGrpcRequestError  = violin + grpc + requestError
	ViolinGrpcGetNodesError = violin + grpc + getNodesError

	ViolinSQLOperationFail = violin + sql + operationFail
	ViolinSQLNoResult      = violin + sql + noResult
	ViolinSQLArgumentError = violin + sql + argumentError
)
