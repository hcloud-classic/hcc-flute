package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"hcc/flute/lib/logger"
	"hcc/flute/model"
)

// ReturnNodes : Publish 'return_nodes' queues to RabbitMQ channel
func ReturnNodes(nodes []model.Node) error {
	qCreate, err := Channel.QueueDeclare(
		"return_nodes",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("return_nodes: Failed to declare a create queue")
		return err
	}

	body, _ := json.Marshal(nodes)
	err = Channel.Publish(
		"",
		qCreate.Name,
		false,
		false,
		amqp.Publishing {
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            body,
		})
	if err != nil {
		logger.Logger.Println("return_nodes: Failed to register publisher")
		return err
	}

	return nil
}

// CreateDHCPDConfig : Publish 'create_dhcpd_config' queues to RabbitMQ channel
func CreateDHCPDConfig(nodes []model.Node) error {
	qCreate, err := Channel.QueueDeclare(
		"create_dhcpd_config",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("create_dhcpd_config: Failed to declare a create queue")
		return err
	}

	body, _ := json.Marshal(nodes)
	err = Channel.Publish(
		"",
		qCreate.Name,
		false,
		false,
		amqp.Publishing {
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            body,
		})
	if err != nil {
		logger.Logger.Println("create_dhcpd_config: Failed to register publisher")
		return err
	}

	return nil
}


