package errors

const (
	PiccoloInternalInitFail            = piccolo + internal + initFail
	PiccoloInternalConnectionFail      = piccolo + internal + connectionFail
	PiccoloInternalUUIDGenerationError = piccolo + internal + UUIDGenerationError

	PiccoloGrpcRequestError = piccolo + grpc + requestError

	PiccoloGraphQLTimestampConversionError = piccolo + graphql + timestampConversionError
	PiccoloGraphQLArgumentError            = piccolo + graphql + argumentError
	PiccoloGraphQLLoginFailed              = piccolo + graphql + loginFailed
	PiccoloGraphQLTokenGenerationError     = piccolo + graphql + tokenGenerationError
	PiccoloGraphQLInvalidToken             = piccolo + graphql + invalidToken
	PiccoloGraphQLTokenExpired             = piccolo + graphql + tokenExpired

	PiccoloMySQLPrepareError = piccolo + sql + prepareError
	PiccoloMySQLExecuteError = piccolo + sql + executeError
)
