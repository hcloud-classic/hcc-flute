package config

type grpc struct {
	Port int64 `goconf:"http:port"` // Port : Port number for listening gRPC request
}

// Grpc : Grpc config structure
var Grpc grpc
