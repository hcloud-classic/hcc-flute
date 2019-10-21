package rabbitmq

import (
	"encoding/json"
	"hcc/flute/dao"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
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

			var node model.Node
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
				logger.Logger.Println("on_node: UUID = " + uuid + ": " + err.Error())
				logger.Logger.Println("on_node: UUID = " + uuid + ": failed to get bmc IP of the node")
				return
			}

			serialNo, err := ipmi.GetSerialNo(bmcIP)
			if err != nil {
				logger.Logger.Println("on_node: UUID = " + uuid + ": " + err.Error())
				logger.Logger.Println("on_node: UUID = " + uuid + ": failed to get serial no of the node")
				return
			}

			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "On" {
				logger.Logger.Println("on_node: UUID = " + uuid + ": already turned on")
				return
			}

			result, err := ipmi.ChangePowerState(bmcIP, serialNo, "On")
			if err != nil {
				logger.Logger.Println("on_node: UUID = " + uuid + ": " + err.Error())
				logger.Logger.Println("on_node: UUID = " + uuid + ": failed to turn on the node")
				return
			}
			logger.Logger.Println("on_node: UUID = " + uuid + ": " + result)
		}
	}()

	return nil
}

// OffNode : Consume 'off_node' queues from RabbitMQ channel
func OffNode() error {
	qCreate, err := Channel.QueueDeclare(
		"off_node",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("off_node: Failed to declare a create queue")
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
		logger.Logger.Println("off_node: Failed to register a create consumer")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("off_node: Received a create message: %s", d.Body)

			var node model.Node
			err = json.Unmarshal(d.Body, &node)
			if err != nil {
				logger.Logger.Println("off_node: Failed to unmarshal node data")
				return
			}

			uuid := node.UUID
			forceOff := node.ForceOff

			var bmcIP string

			sql := "select bmc_ip from node where uuid = ?"
			err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
			if err != nil {
				logger.Logger.Println("off_node: UUID = " + uuid + ": " + err.Error())
				logger.Logger.Println("off_node: UUID = " + uuid + ": failed to get bmc IP of the node")
				return
			}

			serialNo, err := ipmi.GetSerialNo(bmcIP)
			if err != nil {
				logger.Logger.Println("off_node: UUID = " + uuid + ": " + err.Error())
				logger.Logger.Println("off_node: UUID = " + uuid + ": failed to get serial no of the node")
				return
			}

			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "Off" {
				logger.Logger.Println("off_node: UUID = " + uuid + ": already turned off")
				return
			}

			changeState := "GracefulShutdown"
			if forceOff {
				changeState = "ForceOff"
			}
			result, err := ipmi.ChangePowerState(bmcIP, serialNo, changeState)
			if err != nil {
				logger.Logger.Println("off_node: UUID = " + uuid + ": " + err.Error())
				logger.Logger.Println("off_node: UUID = " + uuid + ": failed to turn off the node")
				return
			}
			logger.Logger.Println("off_node: UUID = " + uuid + ": " + result)
		}
	}()

	return nil
}

// GetNodes : Consume 'get_nodes' queues from RabbitMQ channel
func GetNodes() error {
	qCreate, err := Channel.QueueDeclare(
		"get_nodes",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("get_nodes: Failed to declare a create queue")
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
		logger.Logger.Println("get_nodes: Failed to register consumer")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("get_nodes: Received a create message: %s", d.Body)

			var server model.Server
			err = json.Unmarshal(d.Body, &server)
			if err != nil {
				logger.Logger.Println("get_nodes: Failed to unmarshal subnet data")
				return
			}

			serverUUID := server.UUID
			nodeNr := server.NodeNr

			nodes, err := dao.GetAvailableNodes()
			if err != nil {
				logger.Logger.Println(err)
				return
			}

			if nodeNr > len(nodes) {
				logger.Logger.Println("get_nodes: Requested nodeNr is lager than available nodes count")
				return
			}

			for i, node := range nodes {
				if i > nodeNr {
					break
				}
				err := dao.UpdateNodeServerUUID(node, serverUUID)
				if err != nil {
					logger.Logger.Println("get_nodes: error occurred while updating server_uuid of node (UUID = " + node.UUID)
					return
				}
			}

			nodesSelected, err := dao.GetNodesOfServer(serverUUID)
			if err != nil {
				logger.Logger.Println(err)
				return
			}

			err = ReturnNodes(nodesSelected)
			if err != nil {
				logger.Logger.Println(err)
				return
			}

			/*
			TODO
			- 1. select * from node where server_uuid is not null
			- 2. 필요한 갯수 만큼 get
			- 3. get 한 노드들에 대해 update node set server_uuid = [server_uuid]
			- 4. return nodeUUIDs: select * from node where server_uuid = [server_uuid]
			5. publish to harp: create_dhcpd_conf
			 */

			//logger.Logger.Println("create_dhcpd_config: UUID = " + uuid + ": " + result)
		}
	}()

	return nil
}