package dao

import (
	dbsql "database/sql"
	"errors"
	gouuid "github.com/nu7hatch/gouuid"
	pb "hcc/flute/action/grpc/rpcflute"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
	"strconv"
	"time"
)

var nodeSelectColumns = "uuid, server_uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, active, created_at"

// ReadNode : Get all of infos of a node by UUID from database.
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

	sql := "select " + nodeSelectColumns + " from node where uuid = ? and available = 1"
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

// ReadNodeList : Get selected infos of nodes from database.
func ReadNodeList(args map[string]interface{}) (interface{}, error) {
	var nodes []model.Node
	var uuid string
	var createdAt time.Time
	var noLimit bool

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

	if !rowOk && !pageOk {
		noLimit = true
	} else if rowOk && pageOk {
		noLimit = false
	} else {
		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
	}

	sql := "select " + nodeSelectColumns + " from node where available = 1"

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

	if !noLimit {
		sql += " order by created_at desc limit ? offset ?"
	}

	logger.Logger.Println("list_node sql : ", sql)

	var stmt *dbsql.Rows
	var err error

	if noLimit {
		stmt, err = mysql.Db.Query(sql)
	} else {
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	}

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

// ReadNodeAll : Get all of infos of nodes from database.
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
	active, activeOk := args["active"].(int)
	var sql string
	var stmt *dbsql.Rows
	var err error

	if !rowOk && !pageOk {
		sql = "select " + nodeSelectColumns + " from node order by created_at desc"
		if activeOk {
			sql = "select " + nodeSelectColumns + " from node where available = 1 and active = " + strconv.Itoa(active) + " order by created_at desc"
		}
		stmt, err = mysql.Db.Query(sql)
	} else if rowOk && pageOk {
		sql = "select " + nodeSelectColumns + " from node order by created_at desc limit ? offset ?"
		if activeOk {
			sql = "select " + nodeSelectColumns + " from node where available = 1 and active = " + strconv.Itoa(active) + " order by created_at desc limit ? offset ?"
		}
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	} else {
		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
	}

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

// ReadNodeNum : Get count of nodes from database.
func ReadNodeNum(args map[string]interface{}) (interface{}, error) {
	var nodeNum model.NodeNum
	var nodeNr int

	sql := "select count(*) from node where available = 1"
	err := mysql.Db.QueryRow(sql).Scan(&nodeNr)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	logger.Logger.Println("Count: ", nodeNr)
	nodeNum.Number = nodeNr

	return nodeNum, nil
}

// CreateNode : Add a node to database.
func CreateNode(args map[string]interface{}) (interface{}, error) {
	out, err := gouuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	uuid := out.String()

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

	sql := "insert into node(uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, active, available, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(node.UUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Description, node.Active)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return node, nil
}

func NodePowerControl(in *pb.ReqNodePowerControl) ([]string, error) {
	nodes := in.GetNodes()
	if nodes == nil {
		return nil, errors.New("need some Nodes")
	}

	var results []string

	var changeState string
	switch in.GetPowerState() {
	case pb.ReqNodePowerControl_ON:
		changeState = "On"
		break
	case pb.ReqNodePowerControl_OFF:
		changeState = "GracefulShutdown"
		break
	case pb.ReqNodePowerControl_FORCE_OFF:
		changeState = "ForceOff"
		break
	case pb.ReqNodePowerControl_FORCE_RESTART:
		changeState = "ForceRestart"
		break
	}

	for _, node := range nodes {
		if len(node.UUID) == 0 {
			continue
		}

		var bmcIP string
		var result string
		var serialNo string

		sql := "select bmc_ip from node where uuid = ?"
		err := mysql.Db.QueryRow(sql, node.UUID).Scan(&bmcIP)
		if err != nil {
			result = err.Error()
			logger.Logger.Println("NodePowerControl(): "+err.Error())
			goto APPEND
		}

		serialNo, err = ipmi.GetSerialNo(bmcIP)
		if err != nil {
			result = "["+bmcIP+"]: "+err.Error()
			logger.Logger.Println("NodePowerControl(): "+result)
			goto APPEND
		}

		if changeState == "On" {
			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "On" {
				result = "["+bmcIP+"]: Already turned on"
				logger.Logger.Println("NodePowerControl(): "+result)
				goto APPEND
			}
		} else if changeState == "GracefulShutdown" ||
			changeState == "ForceOff" {
			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "Off" {
				result = "["+bmcIP+"]: Already turned off"
				logger.Logger.Println("NodePowerControl(): "+result)
				goto APPEND
			}
		}

		result, err = ipmi.ChangePowerState(bmcIP, serialNo, changeState)
		if err != nil {
			result = "["+bmcIP+"]: "+err.Error()
			logger.Logger.Println("NodePowerControl(): "+result)
			goto APPEND
		}
		result = "["+bmcIP+"]: "+result

	APPEND:
		results = append(results, result)
	}

	return results, nil
}

func GetPowerStateNode(args map[string]interface{}) (interface{}, error) {
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

		result, err := ipmi.GetPowerState(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("need a uuid argument")
}

func checkUpdateNodeArgs(args map[string]interface{}) bool {
	_, serverUUIDOk := args["server_uuid"].(string)
	_, bmcMacAddrOk := args["bmc_mac_addr"].(string)
	_, bmcIPOk := args["bmc_ip"].(string)
	_, pxeMacAdrOk := args["pxe_mac_addr"].(string)
	_, statusOk := args["status"].(string)
	_, cpuCoresOk := args["cpu_cores"].(int)
	_, memoryOk := args["memory"].(int)
	_, descriptionOk := args["description"].(string)
	_, activeOk := args["active"].(int)

	return !serverUUIDOk && !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk && !statusOk && !cpuCoresOk && !memoryOk && !descriptionOk && !activeOk
}

// UpdateNode : Update infos of a node.
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
		if checkUpdateNodeArgs(args) {
			return nil, errors.New("need some arguments")
		}

		sql := "update node set"
		var updateSet = ""
		if serverUUIDOk {
			updateSet += " server_uuid = '" + serverUUID + "', "
		}
		if bmcMacAddrOk {
			updateSet += " bmc_mac_addr = '" + bmcMacAddr + "', "
		}
		if bmcIPOk {
			updateSet += " bmc_ip = '" + bmcIP + "', "
		}
		if pxeMacAdrOk {
			updateSet += " pxe_mac_addr = '" + pxeMacAdr + "', "
		}
		if statusOk {
			updateSet += " status = '" + status + "', "
		}
		if cpuCoresOk {
			updateSet += " cpu_cores = " + strconv.Itoa(cpuCores) + ", "
		}
		if memoryOk {
			updateSet += " memory = " + strconv.Itoa(memory) + ", "
		}
		if descriptionOk {
			updateSet += " description = '" + description + "', "
		}
		if activeOk {
			updateSet += " active = " + strconv.Itoa(active) + ", "
		}
		sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

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

// DeleteNode : Delete a node from database.
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
