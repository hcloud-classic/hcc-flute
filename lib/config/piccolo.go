package config

type piccolo struct {
	ServerAddress        string `goconf:"piccolo:piccolo_server_address"`         // ServerAddress : IP address of server which installed piccolo module
	ServerPort           int64  `goconf:"piccolo:piccolo_server_port"`            // ServerPort : Listening port number of piccolo module
	ConnectionTimeOutMs  int64  `goconf:"piccolo:piccolo_connection_timeout_ms"`  // ConnectionTimeOutMs : Timeout for gRPC client connection of piccolo module
	ConnectionRetryCount int64  `goconf:"piccolo:piccolo_connection_retry_count"` // ConnectionRetryCount : Retry count for gRPC client connection of piccolo module
	RequestTimeoutMs     int64  `goconf:"piccolo:piccolo_request_timeout_ms"`     // RequestTimeoutMs : HTTP timeout for GraphQL request to piccolo module
}

// Piccolo : piccolo config structure
var Piccolo piccolo
