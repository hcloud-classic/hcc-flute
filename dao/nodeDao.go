package dao

import (
	"errors"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
	"strconv"
	"time"
)

func CreateNode(args map[string]interface{}) (interface{}, error) {
	bmcIP, bmcIPOk := args["bmc_ip"].(string)
	description, descriptionOk := args["description"].(string)

	if !descriptionOk {
		description = ""
	}

	if bmcIPOk {
		serialNo, err := ipmi.GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		uuid, err := ipmi.GetUUID(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		BMCmac, err := ipmi.GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoBMC), true)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		PXEmac, err := ipmi.GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoPXE), false)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		powerState, err := ipmi.GetPowerState(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		processors, err := ipmi.GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		cpuCores, err := ipmi.GetProcessorsCores(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		memory, err := ipmi.GetTotalSystemMemory(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		node := model.Node{
			UUID:        uuid,
			BmcMacAddr:  BMCmac,
			BmcIP:       bmcIP,
			PXEMacAddr:  PXEmac,
			Status:      powerState,
			CPUCores:    cpuCores,
			Memory:      memory,
			Description: description,
		}

		sql := "insert into node(uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, now())"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(node.UUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Description)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err2
		}
		logger.Logger.Println(result.LastInsertId())

		err = ipmi.BMCIPParserCheckActive(node.BmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		return node, nil
	}

	return nil, errors.New("need bmc_ip argument")
}

func OnNode(args map[string]interface{}) (interface{}, error) {
	uuid, uuidOk := args["uuid"].(string)

	if uuidOk {
		var bmcIP string

		sql := "select bmc_ip from node where uuid = ?"
		err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		serialNo, err := ipmi.GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		state, _ := ipmi.GetPowerState(bmcIP, serialNo)
		if state == "On" {
			return "Already turned on", nil
		}

		result, err := ipmi.ChangePowerState(bmcIP, serialNo, "On")
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("need uuid argument")
}

func OffNode(args map[string]interface{}) (interface{}, error) {
	uuid, uuidOk := args["uuid"].(string)
	forceOff, _ := args["force_off"].(bool)

	if uuidOk {
		var bmcIP string

		sql := "select bmc_ip from node where uuid = ?"
		err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		serialNo, err := ipmi.GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		state, _ := ipmi.GetPowerState(bmcIP, serialNo)
		if state == "Off" {
			return "Already turned off", nil
		}

		changeState := "GracefulShutdown"
		if forceOff {
			changeState = "ForceOff"
		}
		result, err := ipmi.ChangePowerState(bmcIP, serialNo, changeState)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("need uuid argument")
}

func CreateNodeDetail(args map[string]interface{}) (interface{}, error) {
	nodeUUID, nodeUUIDOk := args["node_uuid"].(string)

	if nodeUUIDOk {
		var bmcIP string

		sql := "select bmc_ip from node where uuid = ?"
		err := mysql.Db.QueryRow(sql, nodeUUID).Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		serialNo, err := ipmi.GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		processorModel, err := ipmi.GetProcessorModel(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		processors, err := ipmi.GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		threads, err := ipmi.GetProcessorsThreads(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		nodedetail := model.NodeDetail{
			NodeUUID:      nodeUUID,
			CPUModel:      processorModel,
			CPUProcessors: processors,
			CPUThreads:    threads,
		}

		sql = "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(nodedetail.NodeUUID, nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err2
		}
		logger.Logger.Println(result.LastInsertId())

		return nodedetail, nil
	}

	return nil, errors.New("need node_uuid argument")
}

func UpdateNode(args map[string]interface{}) (interface{}, error) {
	requestUUIDD, requestUUIDDOK := args["uuid"].(string)
	bmcMacAddr, bmcMacAddrOk := args["bmc_mac_addr"].(string)
	bmcIP, bmcIPOk := args["bmc_ip"].(string)
	pxeMacAdr, pxeMacAdrOk := args["pxe_mac_addr"].(string)
	status, statusOk := args["status"].(string)
	cpuCores, cpuCoresOk := args["cpu_cores"].(int)
	memory, memoryOk := args["memory"].(int)
	description, descriptionOk := args["description"].(string)
	active, activeOk := args["active"].(int)

	node := new(model.Node)
	node.UUID = requestUUIDD
	node.BmcMacAddr = bmcMacAddr
	node.BmcIP = bmcIP
	node.PXEMacAddr = pxeMacAdr
	node.Status = status
	node.CPUCores = cpuCores
	node.Memory = memory
	node.Description = description
	node.Active = active

	if requestUUIDDOK {
		if !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk && !statusOk && !cpuCoresOk && !memoryOk && !descriptionOk && !activeOk {
			return nil, nil
		}

		sql := "update node set"
		if bmcMacAddrOk {
			sql += " bmc_mac_addr = '" + bmcMacAddr + "'"
			if bmcIPOk || pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
				sql += ", "
			}
		}
		if bmcIPOk {
			sql += " bmc_ip = '" + bmcIP + "'"
			if pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
				sql += ", "
			}
		}
		if pxeMacAdrOk {
			sql += " pxe_mac_addr = '" + pxeMacAdr + "'"
			if statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
				sql += ", "
			}
		}
		if statusOk {
			sql += " status = '" + status + "'"
			if cpuCoresOk || memoryOk || descriptionOk || activeOk {
				sql += ", "
			}
		}
		if cpuCoresOk {
			sql += " cpu_cores = '" + strconv.Itoa(cpuCores) + "'"
			if memoryOk || descriptionOk || activeOk {
				sql += ", "
			}
		}
		if memoryOk {
			sql += " memory = '" + strconv.Itoa(memory) + "'"
			if descriptionOk || activeOk {
				sql += ", "
			}
		}
		if descriptionOk {
			sql += " description = '" + description + "'"
			if activeOk {
				sql += ", "
			}
		}
		if activeOk {
			sql += " active = '" + strconv.Itoa(active) + "'"
		}
		sql += " where uuid = ?"

		logger.Logger.Println("update_node sql : ", sql)
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, nil
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(node.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, nil
		}
		logger.Logger.Println(result.LastInsertId())
		return node, nil
	}
	return nil, nil
}

func SelectNode(args map[string]interface{}) (interface{}, error) {
	requestedUUID, ok := args["uuid"].(string)
	if ok {
		node := new(model.Node)

		var uuid string
		var BMCmacAddr string
		var bmcIP string
		var pxeMacAddr string
		var status string
		var cpuCores int
		var memory int
		var description string
		var createdAt time.Time
		var active int

		sql := "select * from node where uuid = ?"
		err := mysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid, &BMCmacAddr, &bmcIP, &pxeMacAddr, &status, &cpuCores, &memory, &description, &createdAt, &active)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		node.UUID = uuid
		node.BmcMacAddr = BMCmacAddr
		node.BmcIP = bmcIP
		node.PXEMacAddr = pxeMacAddr
		node.Status = status
		node.CPUCores = cpuCores
		node.Memory = memory
		node.Description = description
		node.CreatedAt = createdAt
		node.Active = active

		return node, nil
	}
	return nil, errors.New("need uuid argument")
}

func ListNode(args map[string]interface{}) (interface{}, error) {
	var nodes []model.Node
	var uuid string
	var createdAt time.Time

	bmcMacAddr, bmcMacAddrOk := args["bmc_mac_addr"].(string)
	bmcIP, bmcIPOk := args["bmc_ip"].(string)
	pxeMacAdr, pxeMacAdrOk := args["pxe_mac_addr"].(string)
	status, statusOk := args["status"].(string)
	cpuCores, cpuCoresOk := args["cpu_cores"].(int)
	memory, memoryOk := args["memory"].(int)
	description, descriptionOk := args["description"].(string)
	active, activeOk := args["active"].(int)
	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	if !rowOk || !pageOk {
		return nil, nil
	}

	sql := "select * from node where"
	if bmcMacAddrOk {
		sql += " bmc_mac_addr = '" + bmcMacAddr + "'"
		if bmcIPOk || pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
			sql += " and"
		}
	}
	if bmcIPOk {
		sql += " bmc_ip = '" + bmcIP + "'"
		if pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
			sql += " and"
		}
	}
	if pxeMacAdrOk {
		sql += " pxe_mac_addr = '" + pxeMacAdr + "'"
		if statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
			sql += " and"
		}
	}
	if statusOk {
		sql += " status = '" + status + "'"
		if cpuCoresOk || memoryOk || descriptionOk || activeOk {
			sql += " and"
		}
	}
	if cpuCoresOk {
		sql += " cpu_cores = '" + strconv.Itoa(cpuCores) + "'"
		if memoryOk || descriptionOk || activeOk {
			sql += " and"
		}
	}
	if memoryOk {
		sql += " memory = '" + strconv.Itoa(memory) + "'"
		if descriptionOk || activeOk {
			sql += " and"
		}
	}
	if descriptionOk {
		sql += " description = '" + description + "'"
		if activeOk {
			sql += " and"
		}
	}
	if activeOk {
		sql += " active = '" + strconv.Itoa(active) + "'"
	}
	sql += " order by created_at desc limit ? offset ?"

	logger.Logger.Println("list_node sql : ", sql)

	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &createdAt, &active)
		if err != nil {
			logger.Logger.Println(err)
		}
		node := model.Node{UUID: uuid, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, CreatedAt: createdAt, Active: active}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func AllNode(args map[string]interface{}) (interface{}, error) {
	var nodes []model.Node
	var uuid string
	var bmcMacAddr string
	var bmcIP string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var description string
	var createdAt time.Time
	var active int
	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	if !rowOk || !pageOk {
		return nil, nil
	}

	sql := "select * from node order by created_at desc limit ? offset ?"
	logger.Logger.Println("list_server sql  : ", sql)
	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &createdAt, &active)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		node := model.Node{UUID: uuid, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, CreatedAt: createdAt, Active: active}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func DetailNode(args map[string]interface{}) (interface{}, error) {
	nodeDetail := new(model.NodeDetail)
	var nodeUUID string
	var cpuModel string
	var cpuProcessors int
	var cpuThreads int
	requestedNodeUUID, requestedNodeUUIDok := args["node_uuid"].(string)
	if requestedNodeUUIDok {
		sql := "select * from node_detail where node_uuid = ?"
		err := mysql.Db.QueryRow(sql, requestedNodeUUID).Scan(&nodeUUID, &cpuModel, &cpuProcessors, &cpuThreads)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		nodeDetail.NodeUUID = nodeUUID
		nodeDetail.CPUModel = cpuModel
		nodeDetail.CPUProcessors = cpuProcessors
		nodeDetail.CPUThreads = cpuThreads

		return nodeDetail, nil
	}
	return nil, errors.New("need node_uuid argument")
}

func NumNode(args map[string]interface{}) (interface{}, error) {
	var nodeNum model.NodeNum
	var nodeNr int

	sql := "select count(*) from node"
	err := mysql.Db.QueryRow(sql).Scan(&nodeNr)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	logger.Logger.Println("Count: ", nodeNr)
	nodeNum.Number = nodeNr

	return nodeNum, nil
}

func GetAvailableNodes() ([]model.Node, error) {
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

func UpdateNodeServerUUID(node model.Node, serverUUID string) error {
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

func GetNodesOfServer(serverUUID string) ([]model.Node, error) {
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