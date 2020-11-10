package config

type violin struct {
	ServerAddress        string `goconf:"violin:violin_server_address"`         // ServerAddress : IP address of server which installed violin module
	ServerPort           int64  `goconf:"violin:violin_server_port"`            // ServerPort : Listening port number of violin module
	ConnectionTimeOutMs  int64  `goconf:"violin:violin_connection_timeout_ms"`  // ConnectionTimeOutMs : Timeout for gRPC client connection of violin module
	ConnectionRetryCount int64  `goconf:"violin:violin_connection_retry_count"` // ConnectionRetryCount : Retry count for gRPC client connection of violin module
	RequestTimeoutMs     int64  `goconf:"violin:violin_request_timeout_ms"`     // RequestTimeoutMs : HTTP timeout for GraphQL request to violin module
}

// Violin : violin config structure
var Violin violin
