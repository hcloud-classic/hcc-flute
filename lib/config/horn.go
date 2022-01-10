package config

type horn struct {
	ServerAddress        string `goconf:"horn:horn_server_address"`         // ServerAddress : IP address of server which installed horn module
	ServerPort           int64  `goconf:"horn:horn_server_port"`            // ServerPort : Listening port number of horn module
	ConnectionTimeOutMs  int64  `goconf:"horn:horn_connection_timeout_ms"`  // ConnectionTimeOutMs : Timeout for gRPC client connection of horn module
	ConnectionRetryCount int64  `goconf:"horn:horn_connection_retry_count"` // ConnectionRetryCount : Retry count for gRPC client connection of horn module
	RequestTimeoutMs     int64  `goconf:"horn:horn_request_timeout_ms"`     // RequestTimeoutMs : HTTP timeout for GraphQL request to horn module
}

// Horn : horn config structure
var Horn horn
