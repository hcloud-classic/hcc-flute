package dao

import (
	"errors"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
	"strconv"
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