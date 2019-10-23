package dao

import (
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/uuidgen"
	"hcc/flute/model"
	"strconv"
	"time"
)

// ReadNode - cgs
func ReadNode(args map[string]interface{}) (interface{}, error) {
	var node model.Node
	var err error

	uuid := args["uuid"].(string)
	var serverUUID string
	var bmcMacAddr string
	var bmcIP string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var description string
	var createdAt time.Time
	var active int

	sql := "select * from node where uuid = ?"
	err = mysql.Db.QueryRow(sql, uuid).Scan(
		&uuid,
		&serverUUID,
		&bmcMacAddr,
		&bmcIP,
		&pxeMacAdr,
		&status,
		&cpuCores,
		&memory,
		&description,
		&active,
		&createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	node.UUID = uuid
	node.ServerUUID = serverUUID
	node.BmcMacAddr = bmcMacAddr
	node.BmcIP = bmcIP
	node.PXEMacAddr = pxeMacAdr
	node.Status = status
	node.CPUCores = cpuCores
	node.Memory = memory
	node.Description = description
	node.Active = active
	node.CreatedAt = createdAt

	return node, nil
}

// ReadNodeList - cgs
func ReadNodeList(args map[string]interface{}) (interface{}, error) {
	var nodes []model.Node
	var uuid string
	var createdAt time.Time

	serverUUID, serverUUIDOk := args["server_uuid"].(string)
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

	sql := "select * from node where 1 = 1 and"
	if serverUUIDOk {
		sql += " server_uuid = '" + serverUUID + "'"
		if bmcMacAddrOk || bmcIPOk || pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
			sql += " and"
		}
	}
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
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &serverUUID, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &createdAt, &active)
		if err != nil {
			logger.Logger.Println(err)
		}
		node := model.Node{UUID: uuid, ServerUUID: serverUUID, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, CreatedAt: createdAt, Active: active}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// ReadNodeAll - cgs
func ReadNodeAll(args map[string]interface{}) (interface{}, error) {
	var nodes []model.Node
	var uuid string
	var serverUUID string
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
	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &serverUUID, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &active, &createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		node := model.Node{UUID: uuid, ServerUUID: serverUUID, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, Active: active, CreatedAt: createdAt}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// ReadNodeNum - cgs
func ReadNodeNum(args map[string]interface{}) (interface{}, error) {
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

// CreateNode - cgs
func CreateNode(args map[string]interface{}) (interface{}, error) {
	uuid, err := uuidgen.UUIDgen()
	if err != nil {
		logger.Logger.Println("Failed to generate uuid!")
		return nil, err
	}

	node := model.Node{
		UUID:        uuid,
		BmcMacAddr:  args["bmc_mac_addr"].(string),
		BmcIP:       args["bmc_ip"].(string),
		PXEMacAddr:  args["pxe_mac_addr"].(string),
		Status:      args["status"].(string),
		CPUCores:    args["cpu_cores"].(int),
		Memory:      args["memory"].(int),
		Description: args["description"].(string),
		Active:      args["active"].(int),
	}

	sql := "insert into node(uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, active, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(node.UUID, node.ServerUUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Description, node.Active)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return node, nil
}

// UpdateNode - cgs
func UpdateNode(args map[string]interface{}) (interface{}, error) {
	requestUUIDD, requestUUIDDOK := args["uuid"].(string)
	serverUUID, serverUUIDOk := args["server_uuid"].(string)
	bmcMacAddr, bmcMacAddrOk := args["bmc_mac_addr"].(string)
	bmcIP, bmcIPOk := args["bmc_ip"].(string)
	pxeMacAdr, pxeMacAdrOk := args["pxe_mac_addr"].(string)
	status, statusOk := args["status"].(string)
	cpuCores, cpuCoresOk := args["cpu_cores"].(int)
	memory, memoryOk := args["memory"].(int)
	description, descriptionOk := args["description"].(string)
	active, activeOk := args["active"].(int)

	node := new(model.Node)
	node.ServerUUID = serverUUID
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
		if !serverUUIDOk && !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk && !statusOk && !cpuCoresOk && !memoryOk && !descriptionOk && !activeOk {
			return nil, nil
		}

		sql := "update node set"
		if serverUUIDOk {
			sql += " server_uuid = '" + serverUUID + "'"
			if bmcMacAddrOk || bmcIPOk || pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
				sql += ", "
			}
		}
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
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(node.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err2
		}
		logger.Logger.Println(result.LastInsertId())
		return node, nil
	}
	return nil, nil
}

// DeleteNode - cgs
func DeleteNode(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, ok := args["uuid"].(string)
	if ok {
		sql := "delete from node where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(requestedUUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.RowsAffected())

		return requestedUUID, nil
	}

	return requestedUUID, err
}
