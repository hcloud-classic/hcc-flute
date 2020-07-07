package dao

import (
	"errors"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
	"strconv"
	"time"
)

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

func checkReadNodeListPageRow(args map[string]interface{}) bool {
	_, rowOk := args["row"].(int)
	_, pageOk := args["page"].(int)

	return !rowOk || !pageOk
}

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
	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkReadNodeListPageRow(args) {
		return nil, errors.New("need row and page arguments")
	}

	sql := "select * from node where 1=1"

	if serverUUIDOk {
		sql += " and server_uuid = '" + serverUUID + "'"
	}
	if bmcMacAddrOk {
		sql += " and bmc_mac_addr = '" + bmcMacAddr + "'"
	}
	if bmcIPOk {
		sql += " and bmc_ip = '" + bmcIP + "'"
	}
	if pxeMacAdrOk {
		sql += " and pxe_mac_addr = '" + pxeMacAdr + "'"
	}
	if statusOk {
		sql += " and status = '" + status + "'"
	}
	if cpuCoresOk {
		sql += " and cpu_cores = " + strconv.Itoa(cpuCores)
	}
	if memoryOk {
		sql += " and memory = " + strconv.Itoa(memory)
	}
	if descriptionOk {
		sql += " and description = '" + description + "'"
	}
	if activeOk {
		sql += " and active = " + strconv.Itoa(active)
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
		err := stmt.Scan(&uuid, &serverUUID, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &active, &createdAt)
		if err != nil {
			logger.Logger.Println(err)
		}
		node := model.Node{UUID: uuid, ServerUUID: serverUUID, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, Active: active, CreatedAt: createdAt}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
