package rabbitmq

import (
	"encoding/json"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
	"log"
)

func getAvailableNodes() ([]model.Node, error) {
	var nodes []model.Node
	var node model.Node

	sql := "select * from node where server_uuid is not null"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&node.UUID, &node.BmcMacAddr, &node.BmcIP, &node.PXEMacAddr, &node.Status, &node.CPUCores, &node.Memory, &node.Description, &node.CreatedAt, &node.Active)
		if err != nil {
			logger.Logger.Println(err)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func updateNodeServerUUID(node model.Node, serverUUID string) error {
	sql := "update node set server_uuid = server_uuid where uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err2 := stmt.Exec(node.UUID)
	if err2 != nil {
		return err2
	}

	return nil
}

func getNodesOfServer(serverUUID string) ([]model.Node, error) {
	var nodes []model.Node
	var node model.Node

	sql := "select * from node where server_uuid  = " + serverUUID
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&node.UUID, &node.BmcMacAddr, &node.BmcIP, &node.PXEMacAddr, &node.Status, &node.CPUCores, &node.Memory, &node.Description, &node.CreatedAt, &node.Active)
		if err != nil {
			logger.Logger.Println(err)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

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

			nodes, err := getAvailableNodes()
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
				err := updateNodeServerUUID(node, serverUUID)
				if err != nil {
					logger.Logger.Println("get_nodes: error occurred while updating server_uuid of node (UUID = " + node.UUID)
					return
				}
			}

			nodesSelected, err := getNodesOfServer(serverUUID)
			if err != nil {
				logger.Logger.Println(err)
				return
			}

			logger.Logger.Println("get_nodes: publishing ReturnNodes for serverUUID = " + serverUUID)

			err = ReturnNodes(nodesSelected)
			if err != nil {
				logger.Logger.Println(err)
				return
			}
		}
	}()

	return nil
}