package dao

import (
	dbsql "database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	gouuid "github.com/nu7hatch/gouuid"
	pb "hcc/flute/action/grpc/rpcflute"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"strconv"
	"time"
)

var nodeSelectColumns = "uuid, server_uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, active, created_at"

// ReadNode : Get all of infos of a node by UUID from database.
func ReadNode(uuid string) (*pb.Node, error) {
	var node pb.Node

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
	err := mysql.Db.QueryRow(sql, uuid).Scan(
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
	node.CPUCores = int32(cpuCores)
	node.Memory = int32(memory)
	node.Description = description
	node.Active = int32(active)

	node.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	return &node, nil
}

// ReadNodeList : Get selected infos of nodes from database.
func ReadNodeList(in *pb.ReqGetNodeList) (*pb.ResGetNodeList, error) {
	var nodeList pb.ResGetNodeList
	var nodes []pb.Node
	var pnodes []*pb.Node

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

	var isLimit bool
	row := in.GetRow()
	rowOk := row != 0
	page := in.GetPage()
	pageOk := page != 0
	if !rowOk && !pageOk {
		isLimit = false
	} else if rowOk && pageOk {
		isLimit = true
	} else {
		return nil, errors.New("please insert row and page arguments or leave arguments as empty state")
	}

	sql := "select " + nodeSelectColumns + " from node where available = 1"

	if in.Node != nil {
		reqNode := in.Node

		serverUUID = reqNode.ServerUUID
		serverUUIDOk := len(serverUUID) != 0
		bmcMacAddr = reqNode.BmcMacAddr
		bmcMacAddrOk := len(bmcMacAddr) != 0
		bmcIP = reqNode.BmcIP
		bmcIPOk := len(bmcIP) != 0
		pxeMacAdr = reqNode.PXEMacAddr
		pxeMacAdrOk := len(pxeMacAdr) != 0
		status = reqNode.Status
		statusOk := len(status) != 0
		cpuCores = int(reqNode.CPUCores)
		cpuCoresOk := cpuCores != 0
		memory = int(reqNode.Memory)
		memoryOk := memory != 0
		description = reqNode.Description
		descriptionOk := len(description) != 0
		active = int(reqNode.Active)
		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		activeOk := active != 0

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
	}

	var stmt *dbsql.Rows
	var err error
	if isLimit {
		sql += " order by created_at desc limit ? offset ?"
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	} else {
		sql += " order by created_at desc"
		stmt, err = mysql.Db.Query(sql)
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

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}

		nodes = append(nodes, pb.Node{
			UUID:        uuid,
			ServerUUID:  serverUUID,
			BmcMacAddr:  bmcMacAddr,
			BmcIP:       bmcIP,
			PXEMacAddr:  pxeMacAdr,
			Status:      status,
			CPUCores:    int32(cpuCores),
			Memory:      int32(memory),
			Description: description,
			Active:      int32(active),
			CreatedAt:   _createdAt,
		})
	}

	for i := range nodes {
		pnodes = append(pnodes, &nodes[i])
	}

	nodeList.Node = pnodes

	return &nodeList, nil
}

// ReadNodeNum : Get count of nodes from database.
func ReadNodeNum() (*pb.ResGetNodeNum, error) {
	var resNodeNum pb.ResGetNodeNum
	var nodeNr int64

	sql := "select count(*) from node where available = 1"
	err := mysql.Db.QueryRow(sql).Scan(&nodeNr)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	resNodeNum.Num = nodeNr

	return &resNodeNum, nil
}

func checkCreateNodeArgs(reqNode *pb.Node) bool {
	serverUUIDOk := len(reqNode.ServerUUID) != 0
	bmcMacAddrOk := len(reqNode.BmcMacAddr) != 0
	bmcIPOk := len(reqNode.BmcIP) != 0
	pxeMacAdrOk := len(reqNode.PXEMacAddr) != 0
	statusOk := len(reqNode.Status) != 0
	cpuCoresOk := reqNode.CPUCores != 0
	memoryOk := reqNode.Memory != 0
	descriptionOk := len(reqNode.Description) != 0

	return !(serverUUIDOk && bmcMacAddrOk && bmcIPOk && pxeMacAdrOk && statusOk && cpuCoresOk && memoryOk && descriptionOk)
}

// CreateNode : Add a node to database.
func CreateNode(in *pb.ReqCreateNode) (*pb.Node, error) {
	reqNode := in.GetNode()
	if reqNode == nil {
		return nil, errors.New("node is nil")
	}

	out, err := gouuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	uuid := out.String()

	if checkCreateNodeArgs(reqNode) {
		return nil, errors.New("some of arguments are missing")
	}

	node := pb.Node{
		UUID:        uuid,
		BmcMacAddr:  reqNode.BmcMacAddr,
		BmcIP:       reqNode.BmcIP,
		PXEMacAddr:  reqNode.PXEMacAddr,
		Status:      reqNode.Status,
		CPUCores:    reqNode.CPUCores,
		Memory:      reqNode.Memory,
		Description: reqNode.Description,
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
	result, err := stmt.Exec(node.UUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Description)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return &node, nil
}

// NodePowerControl : Change power state of nodes
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
			logger.Logger.Println("NodePowerControl(): " + err.Error())
			goto APPEND
		}

		serialNo, err = ipmi.GetSerialNo(bmcIP)
		if err != nil {
			result = "[" + bmcIP + "]: " + err.Error()
			logger.Logger.Println("NodePowerControl(): " + result)
			goto APPEND
		}

		if changeState == "On" {
			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "On" {
				result = "[" + bmcIP + "]: Already turned on"
				logger.Logger.Println("NodePowerControl(): " + result)
				goto APPEND
			}
		} else if changeState == "GracefulShutdown" ||
			changeState == "ForceOff" {
			state, _ := ipmi.GetPowerState(bmcIP, serialNo)
			if state == "Off" {
				result = "[" + bmcIP + "]: Already turned off"
				logger.Logger.Println("NodePowerControl(): " + result)
				goto APPEND
			}
		}

		result, err = ipmi.ChangePowerState(bmcIP, serialNo, changeState)
		if err != nil {
			result = "[" + bmcIP + "]: " + err.Error()
			logger.Logger.Println("NodePowerControl(): " + result)
			goto APPEND
		}
		result = "[" + bmcIP + "]: " + result

	APPEND:
		results = append(results, result)
	}

	return results, nil
}

// GetNodePowerState : Get power state of the node
func GetNodePowerState(in *pb.ReqNodePowerState) (string, error) {
	uuid := in.GetUUID()
	uuidOk := len(uuid) != 0
	if !uuidOk {
		return "", errors.New("need a uuid argument")
	}

	var bmcIP string

	sql := "select bmc_ip from node where uuid = ?"
	err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
	if err != nil {
		logger.Logger.Println(err)
		return "", err
	}

	serialNo, err := ipmi.GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println(err)
		return "", err
	}

	result, err := ipmi.GetPowerState(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println(err)
		return "", err
	}

	return result, nil
}

func checkUpdateNodeArgs(reqNode *pb.Node) bool {
	serverUUIDOk := len(reqNode.ServerUUID) != 0
	bmcMacAddrOk := len(reqNode.BmcMacAddr) != 0
	bmcIPOk := len(reqNode.BmcIP) != 0
	pxeMacAdrOk := len(reqNode.PXEMacAddr) != 0
	statusOk := len(reqNode.Status) != 0
	cpuCoresOk := reqNode.CPUCores != 0
	memoryOk := reqNode.Memory != 0
	descriptionOk := len(reqNode.Description) != 0
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := reqNode.Active != 0

	return !serverUUIDOk && !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk && !statusOk && !cpuCoresOk && !memoryOk && !descriptionOk && !activeOk
}

// UpdateNode : Update infos of the node.
func UpdateNode(in *pb.ReqUpdateNode) (*pb.Node, error) {
	if in.Node == nil {
		return nil, errors.New("node is nil")
	}
	reqNode := in.Node

	requestedUUID := reqNode.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, errors.New("need a uuid argument")
	}

	if checkUpdateNodeArgs(reqNode) {
		return nil, errors.New("need some arguments")
	}

	var serverUUID string
	var bmcMacAddr string
	var bmcIP string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var description string
	var active int

	serverUUID = in.GetNode().ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	bmcMacAddr = in.GetNode().BmcMacAddr
	bmcMacAddrOk := len(bmcMacAddr) != 0
	bmcIP = in.GetNode().BmcIP
	bmcIPOk := len(bmcIP) != 0
	pxeMacAdr = in.GetNode().PXEMacAddr
	pxeMacAdrOk := len(pxeMacAdr) != 0
	status = in.GetNode().Status
	statusOk := len(status) != 0
	cpuCores = int(in.GetNode().CPUCores)
	cpuCoresOk := cpuCores != 0
	memory = int(in.GetNode().Memory)
	memoryOk := memory != 0
	description = in.GetNode().Description
	descriptionOk := len(description) != 0
	active = int(in.GetNode().Active)
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := active != 0

	node := new(pb.Node)
	node.ServerUUID = serverUUID
	node.UUID = requestedUUID
	node.BmcMacAddr = bmcMacAddr
	node.BmcIP = bmcIP
	node.PXEMacAddr = pxeMacAdr
	node.Status = status
	node.CPUCores = int32(cpuCores)
	node.Memory = int32(memory)
	node.Description = description
	node.Active = int32(active)

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

// DeleteNode : Delete a node from database.
func DeleteNode(in *pb.ReqDeleteNode) (string, error) {
	var err error

	requestedUUID := in.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return "", errors.New("need a uuid argument")
	}

	sql := "delete from node where uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return "", err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err2 := stmt.Exec(requestedUUID)
	if err2 != nil {
		logger.Logger.Println(err2)
		return "", err
	}
	logger.Logger.Println(result.RowsAffected())

	return requestedUUID, nil
}
