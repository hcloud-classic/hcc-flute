package config

type rabbitmq struct {
	ID       string `goconf:"rabbitmq:rabbitmq_id"`       // ID : RabbitMQ server login id
	Password string `goconf:"rabbitmq:rabbitmq_password"` // Password : RabbitMQ server login password
	Address  string `goconf:"rabbitmq:rabbitmq_address"`  // Address : RabbitMQ server address
	Port     int64  `goconf:"rabbitmq:rabbitmq_port"`     // Port : RabbitMQ server port number
}

// RabbitMQ : rabbitmq config structure
var RabbitMQ rabbitmq
