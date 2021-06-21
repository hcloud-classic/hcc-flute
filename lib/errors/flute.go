package errors

const (
	FluteInternalInitFail                 = flute + internal + initFail
	FluteInternalOperationFail            = flute + internal + operationFail
	FluteInternalUUIDGenerationError      = flute + internal + UUIDGenerationError
	FluteInternalTimeStampConversionError = flute + internal + timestampConversionError
	FluteInternalIPMIError                = flute + internal + ipmiError

	FluteGrpcArgumentError = flute + grpc + argumentError
	FluteGrpcRequestError  = flute + grpc + requestError

	FluteSQLOperationFail = flute + sql + operationFail
	FluteSQLNoResult      = flute + sql + noResult
	FluteSQLArgumentError = flute + sql + argumentError
)
