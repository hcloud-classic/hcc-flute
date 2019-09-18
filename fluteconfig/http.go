package fluteconfig

type http struct {
	Port int64 `goconf:"http:port"`  // Port : Port number for listening graphql request via http server
}

var Http http