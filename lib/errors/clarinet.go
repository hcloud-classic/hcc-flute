package errors

const (
	ClarinetInternalInitFail       = clarinet + internal + initFail
	ClarinetInternalConnectionFail = clarinet + internal + connectionFail
	ClarinetInternalParsingError   = clarinet + internal + parsingError

	ClarinetDriverRequestError       = clarinet + driver + requestError
	ClarinetDriverResponseError      = clarinet + driver + responseError
	ClarinetDriverReceiveError       = clarinet + driver + receiveError
	ClarinetDriverParsingError       = clarinet + driver + parsingError
	ClarinetDriverJsonUnmarshalError = clarinet + driver + jsonUnmarshalError

	ClarinetGraphQLArgumentError      = clarinet + graphql + argumentError
	ClarinetGraphQLRequestError       = clarinet + graphql + requestError
	ClarinetGraphQLSendError          = clarinet + graphql + sendError
	ClarinetGraphQLJsonUnmarshalError = clarinet + graphql + jsonUnmarshalError
	ClarinetGraphQLParsingError       = clarinet + graphql + parsingError
)
