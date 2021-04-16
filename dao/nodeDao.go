package dao

import (
	dbsql "database/sql"
	"github.com/golang/protobuf/ptypes"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/iputil"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nodeSelectColumns = "uuid, group_id, server_uuid, node_num, node_ip, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, " +
	"nic_model, nic_speed_mbps, bmc_nic_model, bmc_nic_speed_mbps, description, rack_number, charge_cpu, charge_memory, charge_nic, active, created_at"

// ReadNode : Get all of infos of a node by UUID from database.
func ReadNode(uuid string) (*pb.Node, uint64, string) {
	var node pb.Node

	var groupID int64
	var serverUUID string
	var nodeNum int
	var nodeIP string
	var bmcMacAddr string
	var bmcIPCIDR string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var nicModel string
	var nicSpeedMbps int
	var bmcNicModel string
	var bmcNicSpeedMbps int
	var description string
	var rackNumber int
	var chargeCPU int
	var chargeMemory int
	var chargeNIC int
	var createdAt time.Time
	var active int

	sql := "select " + nodeSelectColumns + " from node where uuid = ? and available = 1"
	row := mysql.Db.QueryRow(sql, uuid)
	err := mysql.QueryRowScan(row,
		&uuid,
		&groupID,
		&serverUUID,
		&nodeNum,
		&nodeIP,
		&bmcMacAddr,
		&bmcIPCIDR,
		&pxeMacAdr,
		&status,
		&cpuCores,
		&memory,
		&nicModel,
		&nicSpeedMbps,
		&bmcNicModel,
		&bmcNicSpeedMbps,
		&description,
		&rackNumber,
		&chargeCPU,
		&chargeMemory,
		&chargeNIC,
		&active,
		&createdAt)
	if err != nil {
		errStr := "ReadNode(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.FluteSQLNoResult, errStr
		}
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	netIP, netIPNet, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		return nil, hcc_errors.FluteGrpcRequestError, "ReadNode(): " + err.Error()
	}
	bmcIP := netIP.String()
	bmcIPSubnetMask := net.IP(netIPNet.Mask).To4().String()

	node.UUID = uuid
	node.GroupID = groupID
	node.ServerUUID = serverUUID
	node.NodeNum = int32(nodeNum)
	node.NodeIP = nodeIP
	node.BmcMacAddr = bmcMacAddr
	node.BmcIP = bmcIP
	node.BmcIPSubnetMask = bmcIPSubnetMask
	node.PXEMacAddr = pxeMacAdr
	node.Status = status
	node.CPUCores = int32(cpuCores)
	node.Memory = int32(memory)
	node.NicModel = nicModel
	node.NicSpeedMbps = int32(nicSpeedMbps)
	node.BmcNicModel = bmcNicModel
	node.BmcNicSpeedMbps = int32(bmcNicSpeedMbps)
	node.Description = description
	node.RackNumber = int32(rackNumber)
	node.ChargeCPU = int32(chargeCPU)
	node.ChargeMemory = int32(chargeMemory)
	node.ChargeNIC = int32(chargeNIC)
	node.Active = int32(active)

	node.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		errStr := "ReadNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteInternalTimeStampConversionError, errStr
	}

	return &node, 0, ""
}

// ReadNodeList : Get selected infos of nodes from database.
func ReadNodeList(in *pb.ReqGetNodeList) (*pb.ResGetNodeList, uint64, string) {
	var nodeList pb.ResGetNodeList
	var nodes []pb.Node
	var pnodes []*pb.Node

	var uuid string
	var groupID int64
	var serverUUID string
	var nodeNum int
	var nodeIP string
	var bmcMacAddr string
	var bmcIPCIDR string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var nicModel string
	var nicSpeedMbps int
	var bmcNicModel string
	var bmcNicSpeedMbps int
	var description string
	var rackNumber int
	var chargeCPU int
	var chargeMemory int
	var chargeNIC int
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
		return nil, hcc_errors.FluteGrpcArgumentError, "ReadNodeList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select " + nodeSelectColumns + " from node where available = 1"

	if in.Node != nil {
		reqNode := in.Node

		uuid = reqNode.UUID
		uuidOk := len(uuid) != 0
		groupID = reqNode.GroupID
		groupIDOk := groupID != 0
		serverUUID = reqNode.ServerUUID
		serverUUIDOk := len(serverUUID) != 0
		nodeNum = int(reqNode.NodeNum)
		// gRPC use 0 value for unset. So I will use -1 for unset node_num. - ish
		nodeNumOk := nodeNum != 0
		nodeIP = reqNode.NodeIP
		nodeIPOk := len(nodeIP) != 0
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
		nicModel = reqNode.NicModel
		nicModelOk := len(nicModel) != 0
		nicSpeedMbps = int(reqNode.NicSpeedMbps)
		nicSpeedMbpsOk := nicSpeedMbps != 0
		bmcNicModel = reqNode.BmcNicModel
		bmcNicModelOk := len(bmcNicModel) != 0
		bmcNicSpeedMbps = int(reqNode.BmcNicSpeedMbps)
		bmcNicSpeedMbpsOk := bmcNicSpeedMbps != 0
		description = reqNode.Description
		descriptionOk := len(description) != 0
		rackNumber = int(reqNode.GetRackNumber())
		rackNumberOk := rackNumber != 0
		chargeCPU = int(reqNode.GetChargeCPU())
		chargeCPUOk := chargeCPU != 0
		chargeMemory = int(reqNode.GetChargeMemory())
		chargeMemoryOk := chargeMemory != 0
		chargeNIC = int(reqNode.GetChargeNIC())
		chargeNICOk := chargeNIC != 0
		active = int(reqNode.Active)
		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		activeOk := active != 0

		if groupIDOk {
			sql += " and group_id = " + strconv.Itoa(int(groupID))
		}
		if uuidOk {
			sql += " and uuid = '" + uuid + "'"
		}
		if serverUUIDOk {
			sql += " and server_uuid = '" + serverUUID + "'"
		}
		if nodeNumOk {
			sql += " and node_num = " + strconv.Itoa(nodeNum)
		}
		if nodeIPOk {
			sql += " and node_ip = '" + nodeIP + "'"
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
		if nicModelOk {
			sql += " and nic_model = '" + nicModel + "'"
		}
		if nicSpeedMbpsOk {
			sql += " and nic_speed_mbps = " + strconv.Itoa(nicSpeedMbps)
		}
		if bmcNicModelOk {
			sql += " and bmc_nic_model = '" + bmcNicModel + "'"
		}
		if bmcNicSpeedMbpsOk {
			sql += " and bmc_nic_speed_mbps = " + strconv.Itoa(bmcNicSpeedMbps)
		}
		if descriptionOk {
			sql += " and description = '" + description + "'"
		}
		if rackNumberOk {
			sql += " and rack_number = " + strconv.Itoa(rackNumber)
		}
		if chargeCPUOk {
			sql += " and charge_cpu = " + strconv.Itoa(chargeCPU)
		}
		if chargeMemoryOk {
			sql += " and charge_memory = " + strconv.Itoa(chargeMemory)
		}
		if chargeNICOk {
			sql += " and charge_nic = " + strconv.Itoa(chargeNIC)
		}
		if activeOk {
			sql += " and active = " + strconv.Itoa(active)
		}
	}

	var stmt *dbsql.Rows
	var err error
	if isLimit {
		sql += " order by created_at desc limit ? offset ?"
		stmt, err = mysql.Query(sql, row, row*(page-1))
	} else {
		sql += " order by created_at desc"
		stmt, err = mysql.Query(sql)
	}

	if err != nil {
		errStr := "ReadNodeList(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &groupID, &serverUUID, &nodeNum, &nodeIP, &bmcMacAddr, &bmcIPCIDR, &pxeMacAdr, &status, &cpuCores, &memory,
			&nicSpeedMbps, &description, &rackNumber, &chargeCPU, &chargeMemory, &chargeNIC, &active, &createdAt)
		if err != nil {
			errStr := "ReadNodeList(): " + err.Error()
			logger.Logger.Println(errStr)
			if strings.Contains(err.Error(), "no rows in result set") {
				return nil, hcc_errors.FluteSQLNoResult, errStr
			}
			return nil, hcc_errors.FluteSQLOperationFail, errStr
		}

		if uuid == "" || pxeMacAdr == "" || cpuCores == 0 || memory == 0 {
			logger.Logger.Println("ReadNodeList(): " + bmcIPCIDR + "'s fields have not yet been filled.")
			continue
		}

		_createdAt, err := ptypes.TimestampProto(createdAt)
		if err != nil {
			errStr := "ReadNodeList(): " + err.Error()
			logger.Logger.Println(errStr)
			return nil, hcc_errors.FluteInternalTimeStampConversionError, errStr
		}

		// gRPC use 0 value for unset. So I will use -1 for unset node_num. - ish
		if nodeNum == -1 {
			nodeNum = 0
		}

		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		if active == 9 {
			active = 0
		}

		netIP, netIPNet, err := net.ParseCIDR(bmcIPCIDR)
		if err != nil {
			return nil, hcc_errors.FluteGrpcRequestError, "ReadNodeList(): " + err.Error()
		}
		bmcIP := netIP.String()
		bmcIPSubnetMask := net.IPv4(netIPNet.Mask[0], netIPNet.Mask[1], netIPNet.Mask[2], netIPNet.Mask[3]).To4().String()

		nodes = append(nodes, pb.Node{
			GroupID:         groupID,
			UUID:            uuid,
			ServerUUID:      serverUUID,
			NodeNum:         int32(nodeNum),
			NodeIP:          nodeIP,
			BmcMacAddr:      bmcMacAddr,
			BmcIP:           bmcIP,
			BmcIPSubnetMask: bmcIPSubnetMask,
			PXEMacAddr:      pxeMacAdr,
			Status:          status,
			CPUCores:        int32(cpuCores),
			Memory:          int32(memory),
			NicModel:        nicModel,
			NicSpeedMbps:    int32(nicSpeedMbps),
			BmcNicModel:     bmcNicModel,
			BmcNicSpeedMbps: int32(bmcNicSpeedMbps),
			Description:     description,
			RackNumber:      int32(rackNumber),
			ChargeCPU:       int32(chargeCPU),
			ChargeMemory:    int32(chargeMemory),
			ChargeNIC:       int32(chargeNIC),
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

// CreateNode : Add a node to database.
func CreateNode(in *pb.ReqCreateNode) (*pb.Node, uint64, string) {
	reqNode := in.GetNode()
	if reqNode == nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): node is nil"
	}

	groupIDOk := reqNode.GroupID != 0
	bmcIPOk := len(reqNode.BmcIP) != 0
	nicSpeedMbpsOk := reqNode.NicSpeedMbps != 0
	descriptionOk := len(reqNode.Description) != 0
	chargeCPUOk := reqNode.ChargeCPU != 0
	chargeMemoryOk := reqNode.ChargeMemory != 0
	chargeNICOk := reqNode.ChargeNIC != 0
	if !groupIDOk || !bmcIPOk || !nicSpeedMbpsOk || !descriptionOk || !chargeCPUOk || !chargeMemoryOk || !chargeNICOk {
		return nil, hcc_errors.FluteGrpcRequestError,
			"CreateNode(): need group_id and bmc_ip, nic_speed_mbps, description, charge_cpu, charge_memory, charge_nic arguments"
	}

	err := iputil.CheckCIDRStr(reqNode.BmcIP)
	if err != nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): " + err.Error()
	}

	_, _, err = net.ParseCIDR(reqNode.BmcIP)
	if err != nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): " + err.Error()
	}

	resGetNodeList, errCode, errText := ReadNodeList(&pb.ReqGetNodeList{
		Node: &pb.Node{
			BmcIP: reqNode.BmcIP,
		},
	})
	if errCode != 0 {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): " + errText
	}
	if len(resGetNodeList.Node) != 0 {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): " + reqNode.BmcIP + " is already registered"
	}

	var pbNode *pb.Node
	var wait sync.WaitGroup
	wait.Add(3)

	uuid, err := ipmi.DoUpdateAllNodes(reqNode.BmcIP, &wait, true, reqNode)
	if err != nil {
		return nil, hcc_errors.FluteInternalIPMIError, "CreateNode(): Error occurred while updating node information"
	}

	go func(uuid string, bmcIP string, wait *sync.WaitGroup) {
		err := ipmi.DoUpdateStatusNodes(uuid, bmcIP, "none", wait)
		if err != nil {
			logger.Logger.Println("CreateNode(): " + err.Error())
		}
	}(uuid, reqNode.BmcIP, &wait)

	wait.Wait()

	pbNode, errCode, errText = ReadNode(uuid)
	if errCode != 0 {
		return nil, errCode, "CreateNode(): ReadNode(): " + errText
	}

	return pbNode, 0, ""
}

// NodePowerControl : Change power state of nodes
func NodePowerControl(in *pb.ReqNodePowerControl) ([]string, uint64, string) {
	nodes := in.GetNode()
	if nodes == nil {
		return nil, hcc_errors.FluteGrpcArgumentError, "NodePowerControl(): need some Nodes"
	}

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

	var results = make([]string, len(nodes))

	var wait sync.WaitGroup
	wait.Add(len(nodes))

	for i, node := range nodes {
		go func(i int, results []string, node *pb.Node) {
			if len(node.UUID) == 0 {
				return
			}

			var bmcIPCIDR string
			var bmcIP string
			var netIP net.IP
			var result string
			var serialNo string

			sql := "select bmc_ip from node where uuid = ?"
			row := mysql.Db.QueryRow(sql, node.UUID)
			err := mysql.QueryRowScan(row, &bmcIPCIDR)
			if err != nil {
				result = err.Error()
				logger.Logger.Println("NodePowerControl(): " + err.Error())
				goto APPEND
			}

			netIP, _, err = net.ParseCIDR(bmcIPCIDR)
			if err != nil {
				result = err.Error()
				logger.Logger.Println("NodePowerControl(): " + err.Error())
				goto APPEND
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
			results[i] = result
			wait.Done()
		}(i, results, node)
	}

	wait.Wait()

	return results, 0, ""
}

// GetNodePowerState : Get power state of the node
func GetNodePowerState(in *pb.ReqNodePowerState) (string, uint64, string) {
	uuid := in.GetUUID()
	uuidOk := len(uuid) != 0
	if !uuidOk {
		return "", hcc_errors.FluteGrpcArgumentError, "GetNodePowerState(): need a uuid argument"
	}

	var bmcIPCIDR string

	sql := "select bmc_ip from node where uuid = ?"
	row := mysql.Db.QueryRow(sql, uuid)
	err := mysql.QueryRowScan(row, &bmcIPCIDR)
	if err != nil {
		errStr := "GetNodePowerState(): " + err.Error()
		logger.Logger.Println(errStr)
		return "", hcc_errors.FluteSQLOperationFail, errStr
	}

	netIP, _, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		return "", hcc_errors.FluteGrpcRequestError, "GetNodePowerState(): " + err.Error()
	}
	bmcIP := netIP.String()

	serialNo, err := ipmi.GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println(err)
		return "", hcc_errors.FluteInternalIPMIError, "GetNodePowerState(): " + err.Error()
	}

	result, err := ipmi.GetPowerState(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println(err)
		return "", hcc_errors.FluteInternalIPMIError, "GetNodePowerState(): " + err.Error()
	}

	return result, 0, ""
}

func checkUpdateNodeArgs(reqNode *pb.Node) bool {
	groupIDOk := reqNode.GroupID != 0
	serverUUIDOk := len(reqNode.ServerUUID) != 0
	// gRPC use 0 value for unset. So I will use -1 for unset node_num. - ish
	nodeNumOk := int(reqNode.NodeNum) != 0
	nodeIPOk := len(reqNode.NodeIP) != 0
	bmcMacAddrOk := len(reqNode.BmcMacAddr) != 0
	bmcIPOk := len(reqNode.BmcIP) != 0
	pxeMacAdrOk := len(reqNode.PXEMacAddr) != 0
	statusOk := len(reqNode.Status) != 0
	cpuCoresOk := reqNode.CPUCores != 0
	memoryOk := reqNode.Memory != 0
	nicModelOk := len(reqNode.NicModel) != 0
	nicSpeedMbpsOk := reqNode.NicSpeedMbps != 0
	bmcNicModelOk := len(reqNode.BmcNicModel) != 0
	bmcNicSpeedMbpsOk := reqNode.BmcNicSpeedMbps != 0
	descriptionOk := len(reqNode.Description) != 0
	rackNumberOk := reqNode.RackNumber != 0
	chargeCPUOk := int(reqNode.GetChargeCPU()) != 0
	chargeMemoryOk := int(reqNode.GetChargeMemory()) != 0
	chargeNICOk := int(reqNode.GetChargeNIC()) != 0
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := reqNode.Active != 0

	return !groupIDOk && !serverUUIDOk && !nodeNumOk && !nodeIPOk && !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk &&
		!statusOk && !cpuCoresOk && !memoryOk &&
		!nicModelOk && !nicSpeedMbpsOk && !bmcNicModelOk && !bmcNicSpeedMbpsOk && !descriptionOk && !rackNumberOk &&
		!chargeCPUOk && !chargeMemoryOk && !chargeNICOk && !activeOk
}

// UpdateNode : Update infos of the node.
func UpdateNode(in *pb.ReqUpdateNode) (*pb.Node, uint64, string) {
	if in.Node == nil {
		return nil, hcc_errors.FluteGrpcArgumentError, "UpdateNode(): node is nil"
	}
	reqNode := in.Node

	requestedUUID := reqNode.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, hcc_errors.FluteGrpcArgumentError, "UpdateNode(): need a uuid argument"
	}

	if checkUpdateNodeArgs(reqNode) {
		return nil, hcc_errors.FluteGrpcArgumentError, "UpdateNode(): need some arguments"
	}

	var groupID int64
	var serverUUID string
	var nodeNum int
	var nodeIP string
	var bmcMacAddr string
	var bmcIP string
	var pxeMacAdr string
	var status string
	var cpuCores int
	var memory int
	var nicModel string
	var nicSpeedMbps int
	var bmcNicModel string
	var bmcNicSpeedMbps int
	var description string
	var rackNumber int
	var chargeCPU int
	var chargeMemory int
	var chargeNIC int
	var active int

	groupID = reqNode.GroupID
	groupIDOk := groupID != 0
	serverUUID = reqNode.ServerUUID
	serverUUIDOk := len(serverUUID) != 0
	nodeNum = int(reqNode.NodeNum)
	// gRPC use 0 value for unset. So I will use -1 for unset node_num. - ish
	nodeNumOk := nodeNum != 0
	nodeIP = reqNode.NodeIP
	nodeIPOk := len(nodeIP) != 0
	bmcMacAddr = reqNode.BmcMacAddr
	bmcMacAddrOk := len(bmcMacAddr) != 0
	bmcIP = reqNode.BmcIP
	bmcIPOk := len(reqNode.BmcIP) != 0
	pxeMacAdr = reqNode.PXEMacAddr
	pxeMacAdrOk := len(pxeMacAdr) != 0
	status = reqNode.Status
	statusOk := len(status) != 0
	cpuCores = int(reqNode.CPUCores)
	cpuCoresOk := cpuCores != 0
	memory = int(reqNode.Memory)
	memoryOk := memory != 0
	nicModel = reqNode.NicModel
	nicModelOk := len(nicModel) != 0
	nicSpeedMbps = int(reqNode.NicSpeedMbps)
	nicSpeedMbpsOk := nicSpeedMbps != 0
	bmcNicModel = reqNode.BmcNicModel
	bmcNicModelOk := len(bmcNicModel) != 0
	bmcNicSpeedMbps = int(reqNode.BmcNicSpeedMbps)
	bmcNicSpeedMbpsOk := bmcNicSpeedMbps != 0
	description = reqNode.Description
	descriptionOk := len(description) != 0
	rackNumber = int(reqNode.GetRackNumber())
	rackNumberOk := rackNumber != 0
	chargeCPU = int(reqNode.GetChargeCPU())
	chargeCPUOk := chargeCPU != 0
	chargeMemory = int(reqNode.GetChargeMemory())
	chargeMemoryOk := chargeMemory != 0
	chargeNIC = int(reqNode.GetChargeNIC())
	chargeNICOk := chargeNIC != 0
	active = int(reqNode.Active)
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := active != 0

	sql := "update node set"
	var updateSet = ""
	if groupIDOk {
		updateSet += " group_id = " + strconv.Itoa(int(groupID)) + ", "
	}
	if serverUUIDOk {
		if serverUUID == "-" {
			serverUUID = ""
		}
		updateSet += " server_uuid = '" + serverUUID + "', "
	}
	if nodeNumOk {
		// gRPC use 0 value for unset. So I will use -1 for unset node_num. - ish
		if nodeNum == 0 && nodeNum < -1 {
			return nil, hcc_errors.FluteGrpcRequestError, "node_num value should be -1 for unset or it should be start from 1"
		}
		updateSet += " node_num = " + strconv.Itoa(nodeNum) + ", "
	}
	if nodeIPOk {
		if nodeIP == "-" {
			nodeIP = ""
		}
		updateSet += " node_ip = '" + nodeIP + "', "
	}
	if bmcMacAddrOk {
		updateSet += " bmc_mac_addr = '" + bmcMacAddr + "', "
	}
	if bmcIPOk {
		err := iputil.CheckCIDRStr(bmcIP)
		if err != nil {
			return nil, hcc_errors.FluteGrpcRequestError, "UpdateNode(): " + err.Error()
		}

		_, _, err = net.ParseCIDR(bmcIP)
		if err != nil {
			return nil, hcc_errors.FluteGrpcRequestError, "UpdateNode(): " + err.Error()
		}

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
	if nicModelOk {
		updateSet += " nic_model = '" + nicModel + "', "
	}
	if nicSpeedMbpsOk {
		updateSet += " nic_speed_mbps = " + strconv.Itoa(nicSpeedMbps) + ", "
	}
	if bmcNicModelOk {
		updateSet += " bmc_nic_model = '" + bmcNicModel + "', "
	}
	if bmcNicSpeedMbpsOk {
		updateSet += " bmc_nic_speed_mbps = " + strconv.Itoa(bmcNicSpeedMbps) + ", "
	}
	if descriptionOk {
		updateSet += " description = '" + description + "', "
	}
	if rackNumberOk {
		updateSet += " rack_number = " + strconv.Itoa(rackNumber) + ", "
	}
	if chargeCPUOk {
		updateSet += " charge_cpu = " + strconv.Itoa(chargeCPU) + ", "
	}
	if chargeMemoryOk {
		updateSet += " charge_memory = " + strconv.Itoa(chargeMemory) + ", "
	}
	if chargeNICOk {
		updateSet += " charge_nic = " + strconv.Itoa(chargeNIC) + ", "
	}
	if activeOk {
		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		if active != 1 && active != 9 {
			return nil, hcc_errors.FluteGrpcRequestError, "active value should be 1 for active or 9 for inactive"
		}
		updateSet += " active = " + strconv.Itoa(active) + ", "
	}
	sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

	logger.Logger.Println("update_node sql : ", sql)

	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "UpdateNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err2 := stmt.Exec(requestedUUID)
	if err2 != nil {
		errStr := "UpdateNode(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	node, errCode, errStr := ReadNode(requestedUUID)
	if errCode != 0 {
		logger.Logger.Println("UpdateNode(): " + errStr)
	}

	return node, 0, ""
}

// DeleteNode : Delete a node from database.
func DeleteNode(in *pb.ReqDeleteNode) (*pb.Node, uint64, string) {
	var err error

	requestedUUID := in.GetUUID()
	requestedUUIDOk := len(requestedUUID) != 0
	if !requestedUUIDOk {
		return nil, hcc_errors.FluteGrpcArgumentError, "DeleteNode(): need a uuid argument"
	}

	node, errCode, errText := ReadNode(requestedUUID)
	if errCode != 0 {
		return nil, hcc_errors.FluteGrpcRequestError, "DeleteNode(): " + errText
	}

	sql := "delete from node where uuid = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		errStr := "DeleteNode(): " + err.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err2 := stmt.Exec(requestedUUID)
	if err2 != nil {
		errStr := "DeleteNode(): " + err2.Error()
		logger.Logger.Println(errStr)
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}

	_, errCode, errStr := DeleteNodeDetail(&pb.ReqDeleteNodeDetail{NodeUUID: requestedUUID})
	if errCode != 0 {
		logger.Logger.Println("DeleteNode(): " + errStr)
	}

	return node, 0, ""
}
