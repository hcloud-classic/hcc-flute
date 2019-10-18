package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
	"hcc/flute/config"
	"hcc/flute/logger"
	"strconv"
)

// Connection : RabbitMQ connection variable
var Connection *amqp.Connection

// Channel : RabbitMQ channel variable
var Channel *amqp.Channel

// PrepareChannel : Connect to RabbitMQ server and create channel.
func PrepareChannel() error {
	Connection, err := amqp.Dial("amqp://" + config.RabbitMQ.ID + ":" + config.RabbitMQ.Password + "@" +
		config.RabbitMQ.Address + ":" + strconv.Itoa(int(config.RabbitMQ.Port)))
	if err != nil {
		return errors.New("failed to connect to RabbitMQ server")
	}
	logger.Logger.Println("Connected to RabbitMQ server")

	Channel, err = Connection.Channel()
	if err != nil {
		return errors.New("failed to open a RabbitMQ channel")
	}
	logger.Logger.Println("Opened RabbitMQ channel.")

	return nil
}
