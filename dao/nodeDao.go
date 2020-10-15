package dao

import (
	dbsql "database/sql"
	"github.com/golang/protobuf/ptypes"
	pb "hcc/flute/action/grpc/pb/rpcflute"
	hccerr "hcc/flute/lib/errors"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/iputil"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"net"
	"strconv"
	"strings"
	"time"
)

var nodeSelectColumns = "uuid, server_uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, rack_number, active, created_at"

// ReadNode : Get all of infos of a node by UUID from database.
func ReadNode(uuid string) (*pb.Node, uint64, string) {
	var node pb.Node

	var serverUUID string
	var bmcMacAddr string
	var bmcIPCIDR string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var description string
	var rackNumber int
	var createdAt time.Time
	var active int

	sql := "select " + nodeSelectColumns + " from node where uuid = ? and available = 1"
	err := mysql.Db.QueryRow(sql, uuid).Scan(
		&uuid,
		&serverUUID,
		&bmcMacAddr,
		&bmcIPCIDR,
		&pxeMacAdr,
		&status,
		&cpuCores,
		&memory,
		&description,
		&rackNumber,
		&active,
		&createdAt)
	if err != nil {
		errStr := "ReadNode(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hccerr.FluteSQLNoResult, errStr
		}
		return nil, hccerr.FluteSQLOperationFail, errStr
	}

	netIP, netIPNet, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		return nil, hccerr.FluteGrpcRequestError, "ReadNode(): " + err.Error()
	}
	bmcIP := netIP.String()
	bmcIPSubnetMask := netIPNet.Mask.String()

	node.UUID = uuid
	node.ServerUUID = serverUUID
	node.BmcMacAddr = bmcMacAddr
	node.BmcIP = bmcIP
	node.BmcIPSubnetMask = bmcIPSubnetMask
	node.PXEMacAddr = pxeMacAdr
	node.Status = status
	node.CPUCores = int32(cpuCores)
	node.Memory = int32(memory)
	node.Description = description
	node.RackNumber = int32(rackNumber)
	node.Active = int32(active)

	node.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		errStr := "ReadNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteInternalTimeStampConversionError, errStr
	}

	return &node, 0, ""
}

// ReadNodeList : Get selected infos of nodes from database.
func ReadNodeList(in *pb.ReqGetNodeList) (*pb.ResGetNodeList, uint64, string) {
	var nodeList pb.ResGetNodeList
	var nodes []pb.Node
	var pnodes []*pb.Node

	var uuid string
	var serverUUID string
	var bmcMacAddr string
	var bmcIPCIDR string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var description string
	var rackNumber int
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
		return nil, hccerr.FluteGrpcArgumentError, "ReadNodeList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select " + nodeSelectColumns + " from node where available = 1"

	if in.Node != nil {
		reqNode := in.Node

		uuid = reqNode.UUID
		uuidOk := len(uuid) != 0
		serverUUID = reqNode.ServerUUID
		serverUUIDOk := len(serverUUID) != 0
		bmcMacAddr = reqNode.BmcMacAddr
		bmcMacAddrOk := len(bmcMacAddr) != 0
		bmcIPCIDR = reqNode.BmcIP
		bmcIPOk := len(bmcIPCIDR) != 0
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
		rackNumberOk := rackNumber != 0
		active = int(reqNode.Active)
		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		activeOk := active != 0

		if uuidOk {
			sql += " and uuid = '" + uuid + "'"
		}
		if serverUUIDOk {
			sql += " and server_uuid = '" + serverUUID + "'"
		}
		if bmcMacAddrOk {
			sql += " and bmc_mac_addr = '" + bmcMacAddr + "'"
		}
		if bmcIPOk {
			sql += " and bmc_ip = '" + bmcIPCIDR + "'"
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
		if rackNumberOk {
			sql += " and rack_number = " + strconv.Itoa(rackNumber)
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
		errStr := "ReadNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &serverUUID, &bmcMacAddr, &bmcIPCIDR, &pxeMacAdr, &status, &cpuCores, &memory, &description, &rackNumber, &active, &createdAt)
		if err != nil {
			errStr := "ReadNodeList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hccerr.FluteSQLNoResult, errStr
			}
			return nil, hccerr.FluteSQLOperationFail, errStr
		}

		if uuid == "" || pxeMacAdr == "" || cpuCores == 0 || memory == 0 {
			logger.Logger.Println("ReadNodeList(): " + bmcIPCIDR + "'s fields have not yet been filled.")
			continue
		}

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			errStr := "ReadNodeList(): " + err.Error()
			logger.Logger.Println(errStr)
			return nil, hccerr.FluteInternalTimeStampConversionError, errStr
		}

		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		if active == 9 {
			active = 0
		}

		netIP, netIPNet, err := net.ParseCIDR(bmcIPCIDR)
		if err != nil {
			return nil, hccerr.FluteGrpcRequestError, "ReadNodeList(): " + err.Error()
		}
		bmcIP := netIP.String()
		bmcIPSubnetMask := net.IPv4(netIPNet.Mask[0], netIPNet.Mask[1], netIPNet.Mask[2], netIPNet.Mask[3]).To4().String()

		nodes = append(nodes, pb.Node{
			UUID:            uuid,
			ServerUUID:      serverUUID,
			BmcMacAddr:      bmcMacAddr,
			BmcIP:           bmcIP,
			BmcIPSubnetMask: bmcIPSubnetMask,
			PXEMacAddr:      pxeMacAdr,
			Status:          status,
			CPUCores:        int32(cpuCores),
			Memory:          int32(memory),
			Description:     description,
			RackNumber:      int32(rackNumber),
			Active:          int32(active),
			CreatedAt:       _createdAt,
		})
	}

	for i := range nodes {
		pnodes = append(pnodes, &nodes[i])
	}

	nodeList.Node = pnodes

	return &nodeList, 0, ""
}

// ReadNodeNum : Get count of nodes from database.
func ReadNodeNum() (*pb.ResGetNodeNum, uint64, string) {
	var resNodeNum pb.ResGetNodeNum
	var nodeNr int64

	sql := "select count(*) from node where available = 1"
	err := mysql.Db.QueryRow(sql).Scan(&nodeNr)
	if err != nil {
		errStr := "ReadNodeNum(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hccerr.FluteSQLNoResult, errStr
		}
		return nil, hccerr.FluteSQLOperationFail, errStr
	}
	resNodeNum.Num = nodeNr

	return &resNodeNum, 0, ""
}

func checkCreateNodeArgs(reqNode *pb.Node) bool {
	UUIDOk := len(reqNode.UUID) != 0
	serverUUIDOk := len(reqNode.ServerUUID) != 0
	bmcMacAddrOk := len(reqNode.BmcMacAddr) != 0
	pxeMacAdrOk := len(reqNode.PXEMacAddr) != 0
	statusOk := len(reqNode.Status) != 0
	cpuCoresOk := reqNode.CPUCores != 0
	memoryOk := reqNode.Memory != 0

	return !(UUIDOk && serverUUIDOk && bmcMacAddrOk && pxeMacAdrOk && statusOk && cpuCoresOk && memoryOk)
}

// CreateNode : Add a node to database.
func CreateNode(in *pb.ReqCreateNode) (*pb.Node, uint64, string) {
	reqNode := in.GetNode()
	if reqNode == nil {
		return nil, hccerr.FluteGrpcRequestError, "CreateNode(): node is nil"
	}

	bmcIPOk := len(reqNode.BmcIP) != 0
	descriptionOk := len(reqNode.Description) != 0
	if !bmcIPOk || !descriptionOk {
		return nil, hccerr.FluteGrpcRequestError, "CreateNode(): need bmcIP and description arguments"
	} else if !bmcIPOk && checkCreateNodeArgs(reqNode) {
		return nil, hccerr.FluteGrpcRequestError, "CreateNode(): some of arguments are missing"
	}

	err := iputil.CheckCIDRStr(reqNode.BmcIP)
	if err != nil {
		return nil, hccerr.FluteGrpcRequestError, "CreateNode(): " + err.Error()
	}

	_, _, err = net.ParseCIDR(reqNode.BmcIP)
	if err != nil {
		return nil, hccerr.FluteGrpcRequestError, "CreateNode(): " + err.Error()
	}

	node := pb.Node{
		UUID:        reqNode.UUID,
		ServerUUID:  "",
		BmcMacAddr:  reqNode.BmcMacAddr,
		BmcIP:       reqNode.BmcIP,
		PXEMacAddr:  reqNode.PXEMacAddr,
		Status:      reqNode.Status,
		CPUCores:    reqNode.CPUCores,
		Memory:      reqNode.Memory,
		Description: reqNode.Description,
		RackNumber:  reqNode.RackNumber,
	}

	sql := "insert into node(uuid, server_uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, rack_number, created_at, available) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), 1)"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "CreateNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(node.UUID, node.ServerUUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Description, node.RackNumber)
	if err != nil {
		errStr := "CreateNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}

	return &node, 0, ""
}

// NodePowerControl : Change power state of nodes
func NodePowerControl(in *pb.ReqNodePowerControl) ([]string, uint64, string) {
	nodes := in.GetNode()
	if nodes == nil {
		return nil, hccerr.FluteGrpcArgumentError, "NodePowerControl(): need some Nodes"
	}

	var results []string

	var changeState string
	switch in.GetPowerState() {
	case pb.PowerState_ON:
		changeState = "On"
		break
	case pb.PowerState_OFF:
		changeState = "GracefulShutdown"
		break
	case pb.PowerState_FORCE_OFF:
		changeState = "ForceOff"
		break
	case pb.PowerState_FORCE_RESTART:
		changeState = "ForceRestart"
		break
	}

	for _, node := range nodes {
		if len(node.UUID) == 0 {
			continue
		}

		var bmcIPCIDR string
		var bmcIP string
		var netIP net.IP
		var result string
		var serialNo string

		sql := "select bmc_ip from node where uuid = ?"
		err := mysql.Db.QueryRow(sql, node.UUID).Scan(&bmcIPCIDR)
		if err != nil {
			result = err.Error()
			logger.Logger.Println("NodePowerControl(): " + err.Error())
			goto APPEND
		}

		netIP, _, err = net.ParseCIDR(bmcIPCIDR)
		if err != nil {
			return nil, hccerr.FluteGrpcRequestError, "NodePowerControl(): " + err.Error()
		}
		bmcIP = netIP.String()

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

	return results, 0, ""
}

// GetNodePowerState : Get power state of the node
func GetNodePowerState(in *pb.ReqNodePowerState) (string, uint64, string) {
	uuid := in.GetUUID()
	uuidOk := len(uuid) != 0
	if !uuidOk {
		return "", hccerr.FluteGrpcArgumentError, "GetNodePowerState(): need a uuid argument"
	}

	var bmcIPCIDR string

	sql := "select bmc_ip from node where uuid = ?"
	err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIPCIDR)
	if err != nil {
		errStr := "GetNodePowerState(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.FluteSQLOperationFail, errStr
	}

	netIP, _, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		return "", hccerr.FluteGrpcRequestError, "GetNodePowerState(): " + err.Error()
	}
	bmcIP := netIP.String()

	serialNo, err := ipmi.GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println(err)
		return "", hccerr.FluteInternalIPMIError, "GetNodePowerState(): " + err.Error()
	}

	result, err := ipmi.GetPowerState(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println(err)
		return "", hccerr.FluteInternalIPMIError, "GetNodePowerState(): " + err.Error()
	}

	return result, 0, ""
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
	rackNumberOk := reqNode.RackNumber != 0
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := reqNode.Active != 0

	return !serverUUIDOk && !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk && !statusOk && !cpuCoresOk && !memoryOk && !descriptionOk && !rackNumberOk && !activeOk
}

// UpdateNode : Update infos of the node.
func UpdateNode(in *pb.ReqUpdateNode) (*pb.Node, uint64, string) {
	if in.Node == nil {
		return nil, hccerr.FluteGrpcArgumentError, "UpdateNode(): node is nil"
	}
	reqNode := in.Node

	requestedUUID := reqNode.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, hccerr.FluteGrpcArgumentError, "UpdateNode(): need a uuid argument"
	}

	if checkUpdateNodeArgs(reqNode) {
		return nil, hccerr.FluteGrpcArgumentError, "UpdateNode(): need some arguments"
	}

	var serverUUID string
	var bmcMacAddr string
	var bmcIP string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var description string
	var rackNumber int
	var active int

	serverUUID = in.GetNode().ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	bmcMacAddr = in.GetNode().BmcMacAddr
	bmcMacAddrOk := len(bmcMacAddr) != 0
	bmcIP = in.GetNode().BmcIP
	bmcIPOk := len(reqNode.BmcIP) != 0
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
	rackNumberOk := rackNumber != 0
	active = int(in.GetNode().Active)
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := active != 0

	err := iputil.CheckCIDRStr(bmcIP)
	if err != nil {
		return nil, hccerr.FluteGrpcRequestError, "UpdateNode(): " + err.Error()
	}

	_, _, err = net.ParseCIDR(bmcIP)
	if err != nil {
		return nil, hccerr.FluteGrpcRequestError, "UpdateNode(): " + err.Error()
	}

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
	node.RackNumber = int32(rackNumber)
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
	if rackNumberOk {
		updateSet += " rack_number = " + strconv.Itoa(rackNumber) + ", "
	}
	if activeOk {
		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		if active != 1 && active != 9 {
			return nil, hccerr.FluteGrpcRequestError, "active value should be 1 for active or 9 for inactive"
		}
		updateSet += " active = " + strconv.Itoa(active) + ", "
	}
	sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

	logger.Logger.Println("update_node sql : ", sql)

	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "UpdateNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err2 := stmt.Exec(node.UUID)
	if err2 != nil {
		errStr := "UpdateNode(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hccerr.FluteSQLOperationFail, errStr
	}

	node, errCode, errStr := ReadNode(node.UUID)
	if errCode != 0 {
		logger.Logger.Println("UpdateNode(): " + errStr)
	}

	return node, 0, ""
}

// DeleteNode : Delete a node from database.
func DeleteNode(in *pb.ReqDeleteNode) (string, uint64, string) {
	var err error

	requestedUUID := in.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return "", hccerr.FluteGrpcArgumentError, "DeleteNode(): need a uuid argument"
	}

	sql := "delete from node where uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		errStr := "DeleteNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err2 := stmt.Exec(requestedUUID)
	if err2 != nil {
		errStr := "DeleteNode(): " + err2.Error()
		logger.Logger.Println(errStr)
		return "", hccerr.FluteSQLOperationFail, errStr
	}
	logger.Logger.Println(result.RowsAffected())

	_, errCode, errStr := DeleteNodeDetail(&pb.ReqDeleteNodeDetail{NodeUUID: requestedUUID})
	if errCode != 0 {
		logger.Logger.Println("DeleteNode(): " + errStr)
	}

	return requestedUUID, 0, ""
}
