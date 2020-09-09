package errors

const (
	HarpInternalInitFail                 = harp + internal + initFail
	HarpInternalOperationFail            = harp + internal + operationFail
	HarpInternalUUIDGenerationError      = harp + internal + UUIDGenerationError
	HarpInternalTimeStampConversionError = harp + internal + timestampConversionError
	HarpInternalInterfaceLookupError     = harp + internal + interfaceAddrLookupError
	HarpInternalPFError                  = harp + internal + pfError
	HarpInternalDHCPDError               = harp + internal + dhcpdError
	HarpInternalFileError                = harp + internal + fileError
	HarpInternalIfconfigError            = harp + internal + ifconfigError
	HarpInternalIPAddressError           = harp + internal + IPAddressError
	HarpInternalSubnetInUseError         = harp + internal + subnetInUseError
	HarpInternalSubnetNotAllocatedError  = harp + internal + subnetNotAllocatedError
	HarpInternalAdaptiveIPAllocatedError = harp + internal + adaptiveIPAllocatedError

	HarpGrpcArgumentError = harp + grpc + argumentError
	HarpGrpcRequestError  = harp + grpc + requestError

	HarpSQLOperationFail = harp + sql + operationFail
	HarpSQLNoResult      = harp + sql + noResult
	HarpSQLArgumentError = harp + sql + argumentError
)
