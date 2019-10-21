package config

type http struct {
	Port int64 `goconf:"http:port"` // Port : Port number for listening graphql request via http server
}

// HTTP : http config structure
var HTTP http
