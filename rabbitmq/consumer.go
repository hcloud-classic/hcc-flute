package rabbitmq

import (
	"encoding/json"
	"hcc/flute/ipmi"
	"hcc/flute/logger"
	"hcc/flute/mysql"
	"hcc/flute/types"
	"log"
)

// OnNode : Consume 'on_node' queues from RabbitMQ channel
func OnNode() error {
	qCreate, err := Channel.QueueDeclare(
		"on_node",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("on_node: Failed to declare a create queue")
		return err
	}

	msgsCreate, err := Channel.Consume(
		qCreate.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Println("on_node: Failed to register a create consumer")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("on_node: Received a create message: %s", d.Body)

			var node types.Node
			err = json.Unmarshal(d.Body, &node)
			if err != nil {
				logger.Logger.Println("on_node: Failed to unmarshal node data")
				return
			}

			uuid := node.UUID

			var bmcIP string

			sql := "select bmc_ip from node where uuid = ?"
			err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
			if err != nil {
				logger.Logger.Println("on_node: UUID = "+uuid+": "+err.Error())
				logger.Logger.Println("on_node: UUID = "+uuid+": failed to get bmc IP of the node")
				return
			}

			serialNo, err := ipmi.GetSerialNo(bmcIP)
			if err != nil {
				logger.Logger.Println("on_node: UUID = "+uuid+": "+err.Error())
				logger.Logger.Println("on_node: UUID = "+uuid+": failed to get serial no of the node")
				return
			}

			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "On" {
				logger.Logger.Println("on_node: UUID = "+uuid+": already turned on")
				return
			}

			result, err := ipmi.ChangePowerState(bmcIP, serialNo, "On")
			if err != nil {
				logger.Logger.Println("on_node: UUID = "+uuid+": "+err.Error())
				logger.Logger.Println("on_node: UUID = "+uuid+": failed to turn on the node")
				return
			}
			logger.Logger.Println("on_node: UUID = "+uuid+": "+result)
		}
	}()

	return nil
}
