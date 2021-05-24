package daoext

import (
	dbsql "database/sql"
	"github.com/golang/protobuf/ptypes"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
	"net"
	"strconv"
	"strings"
	"time"
)

var nodeSelectColumns = "uuid, node_name, group_id, server_uuid, node_num, node_ip, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, " +
	"nic_speed_mbps, description, rack_number, active, created_at"

// ReadNode : Get all of infos of a node by UUID from database.
func ReadNode(uuid string) (*pb.Node, uint64, string) {
	var node pb.Node

	var nodeName string
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
	var nicSpeedMbps int
	var description string
	var rackNumber int
	var createdAt time.Time
	var active int

	sql := "select " + nodeSelectColumns + " from node where uuid = ? and available = 1"
	row := mysql.Db.QueryRow(sql, uuid)
	err := mysql.QueryRowScan(row,
		&uuid,
		&nodeName,
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
		&nicSpeedMbps,
		&description,
		&rackNumber,
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
	node.NodeName = nodeName
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
	node.NicSpeedMbps = int32(nicSpeedMbps)
	node.Description = description
	node.RackNumber = int32(rackNumber)
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
	var nodeName string
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
	var nicSpeedMbps int
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
		return nil, hcc_errors.FluteGrpcArgumentError, "ReadNodeList(): please insert row and page arguments or leave arguments as empty state"
	}

	sql := "select " + nodeSelectColumns + " from node where available = 1"

	if in.Node != nil {
		reqNode := in.Node

		uuid = reqNode.UUID
		uuidOk := len(uuid) != 0
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
		nicSpeedMbps = int(reqNode.NicSpeedMbps)
		nicSpeedMbpsOk := nicSpeedMbps != 0
		description = reqNode.Description
		descriptionOk := len(description) != 0
		rackNumber = int(reqNode.GetRackNumber())
		rackNumberOk := rackNumber != 0
		active = int(reqNode.Active)
		// gRPC use 0 value for unset. So I will use 9 value for inactive. - ish
		activeOk := active != 0

		if uuidOk {
			sql += " and uuid like '%" + uuid + "%'"
		}
		if nodeNameOk {
			sql += " and node_name like '%" + nodeName + "%'"
		}
		if groupIDOk {
			sql += " and group_id = " + strconv.Itoa(int(groupID))
		}
		if serverUUIDOk {
			sql += " and server_uuid like '%" + serverUUID + "%'"
		}
		if nodeNumOk {
			sql += " and node_num = " + strconv.Itoa(nodeNum)
		}
		if nodeIPOk {
			sql += " and node_ip like '%" + nodeIP + "%'"
		}
		if bmcMacAddrOk {
			sql += " and bmc_mac_addr like '%" + bmcMacAddr + "%'"
		}
		if bmcIPOk {
			sql += " and bmc_ip like '%" + bmcIPCIDR + "%'"
		}
		if pxeMacAdrOk {
			sql += " and pxe_mac_addr like '%" + pxeMacAdr + "%'"
		}
		if statusOk {
			sql += " and status like '%" + status + "%'"
		}
		if cpuCoresOk {
			sql += " and cpu_cores = " + strconv.Itoa(cpuCores)
		}
		if memoryOk {
			sql += " and memory = " + strconv.Itoa(memory)
		}
		if nicSpeedMbpsOk {
			sql += " and nic_speed_mbps = " + strconv.Itoa(nicSpeedMbps)
		}
		if descriptionOk {
			sql += " and description like '%" + description + "%'"
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
		err := stmt.Scan(&uuid, &nodeName, &groupID, &serverUUID, &nodeNum, &nodeIP, &bmcMacAddr, &bmcIPCIDR, &pxeMacAdr, &status, &cpuCores, &memory,
			&nicSpeedMbps, &description, &rackNumber, &active, &createdAt)
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
			UUID:            uuid,
			NodeName:        nodeName,
			GroupID:         groupID,
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
			NicSpeedMbps:    int32(nicSpeedMbps),
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
func ReadNodeNum(in *pb.ReqGetNodeNum) (*pb.ResGetNodeNum, uint64, string) {
	var resNodeNum pb.ResGetNodeNum
	var nodeNr int64
	var groupID = in.GetGroupID()

	sql := "select count(*) from node where available = 1"
	if groupID != 0 {
		sql = "select count(*) from node where available = 1 and group_id = " + strconv.Itoa(int(groupID))
	}
	row := mysql.Db.QueryRow(sql)
	err := mysql.QueryRowScan(row, &nodeNr)
	if err != nil {
		errStr := "ReadNodeNum(): " + err.Error()
		logger.Logger.Println(errStr)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, hcc_errors.FluteSQLNoResult, errStr
		}
		return nil, hcc_errors.FluteSQLOperationFail, errStr
	}
	resNodeNum.Num = nodeNr

	return &resNodeNum, 0, ""
}
