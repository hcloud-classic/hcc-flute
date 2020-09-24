package errors

const (
	PianoInternalInitFail                 = piano + internal + initFail
	PianoInternalOperationFail            = piano + internal + operationFail

	PianoGrpcArgumentError = piano + grpc + argumentError

	PianoInfluxDBReadMetricError = piano + influxDB + readMetricError
)
