package dao

import (
	"hcc/flute/daoext"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/iputil"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"net"
	"strconv"
	"sync"
)

// CreateNode : Add a node to database.
func CreateNode(in *pb.ReqCreateNode) (*pb.Node, uint64, string) {
	reqNode := in.GetNode()
	if reqNode == nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): node is nil"
	}

	nodeNameOk := len(reqNode.NodeName) != 0
	groupIDOk := reqNode.GroupID != 0
	bmcIPOk := len(reqNode.BmcIP) != 0
	nicSpeedMbpsOk := reqNode.NicSpeedMbps != 0
	descriptionOk := len(reqNode.Description) != 0

	nicDetailDataOk := len(in.NicDetailData) != 0

	if !nodeNameOk || !groupIDOk || !bmcIPOk || !nicSpeedMbpsOk || !descriptionOk || !nicDetailDataOk {
		return nil, hcc_errors.FluteGrpcRequestError,
			"CreateNode(): need node_name, group_id and bmc_ip, nic_speed_mbps, description, " +
				"nic_detail_data arguments"
	}

	err := iputil.CheckCIDRStr(reqNode.BmcIP)
	if err != nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): " + err.Error()
	}

	_, _, err = net.ParseCIDR(reqNode.BmcIP)
	if err != nil {
		return nil, hcc_errors.FluteGrpcRequestError, "CreateNode(): " + err.Error()
	}

	resGetNodeList, errCode, errText := daoext.ReadNodeList(&pb.ReqGetNodeList{
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
	wait.Add(2)

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

	_, errCode, errText = CreateNodeDetail(&pb.ReqCreateNodeDetail{
		NodeDetail: &pb.NodeDetail{
			NodeUUID:       uuid,
			NodeDetailData: "",
			NicDetailData:  in.NicDetailData,
		},
	})
	if errCode != 0 {
		return nil, errCode, "CreateNode(): CreateNodeDetail(): " + errText
	}

	pbNode, errCode, errText = daoext.ReadNode(uuid)
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
	nodeNameOk := len(reqNode.NodeName) != 0
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
	nicSpeedMbpsOk := reqNode.NicSpeedMbps != 0
	descriptionOk := len(reqNode.Description) != 0
	rackNumberOk := reqNode.RackNumber != 0
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := reqNode.Active != 0

	return !nodeNameOk && !groupIDOk && !serverUUIDOk && !nodeNumOk && !nodeIPOk && !bmcMacAddrOk && !bmcIPOk && !pxeMacAdrOk &&
		!statusOk && !cpuCoresOk && !memoryOk &&
		!nicSpeedMbpsOk && !descriptionOk && !rackNumberOk && !activeOk
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

	var nodeName string
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
	var nicSpeedMbps int
	var description string
	var rackNumber int
	var active int

	nodeName = reqNode.NodeName
	nodeNameOk := len(nodeName) != 0
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
	nicSpeedMbps = int(reqNode.NicSpeedMbps)
	nicSpeedMbpsOk := nicSpeedMbps != 0
	description = reqNode.Description
	descriptionOk := len(description) != 0
	rackNumber = int(reqNode.GetRackNumber())
	rackNumberOk := rackNumber != 0
	active = int(reqNode.Active)
	// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
	activeOk := active != 0

	sql := "update node set"
	var updateSet = ""

	if nodeNameOk {
		updateSet += " node_name = '" + nodeName + "', "
	}
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
	if nicSpeedMbpsOk {
		updateSet += " nic_speed_mbps = " + strconv.Itoa(nicSpeedMbps) + ", "
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

	node, errCode, errStr := daoext.ReadNode(requestedUUID)
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

	node, errCode, errText := daoext.ReadNode(requestedUUID)
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
