package config

type harp struct {
	ServerAddress        string `goconf:"harp:harp_server_address"`         // ServerAddress : IP address of server which installed harp module
	ServerPort           int64  `goconf:"harp:harp_server_port"`            // ServerPort : Listening port number of harp module
	ConnectionTimeOutMs  int64  `goconf:"harp:harp_connection_timeout_ms"`  // ConnectionTimeOutMs : Timeout for gRPC client connection of harp module
	ConnectionRetryCount int64  `goconf:"harp:harp_connection_retry_count"` // ConnectionRetryCount : Retry count for gRPC client connection of harp module
	RequestTimeoutMs     int64  `goconf:"harp:harp_request_timeout_ms"`     // RequestTimeoutMs : HTTP timeout for GraphQL request to harp module
}

// Harp : harp config structure
var Harp harp
