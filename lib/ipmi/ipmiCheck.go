package ipmi

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"hcc/flute/action/grpc/client"
	"hcc/flute/daoext"
	"hcc/flute/lib/config"
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

var checkNodeAllLocked = false
var checkNodeStatusLocked = false
var checkServerStatusLocked = false
var updateNodeDetailLocked = make(map[string]bool)

func delayMillisecond(n time.Duration) {
	time.Sleep(n * time.Millisecond)
}

func checkNodeAllLock() {
	checkNodeAllLocked = true
}

func checkNodeAllUnlock() {
	checkNodeAllLocked = false
}

func checkNodeStatusLock() {
	checkNodeStatusLocked = true
}

func checkNodeStatusUnlock() {
	checkNodeStatusLocked = false
}

func checkServerStatusLock() {
	checkServerStatusLocked = true
}

func checkServerStatusUnlock() {
	checkServerStatusLocked = false
}

func isUpdateNodeDetailLocked(uuid string) bool {
	isLocked, exist := updateNodeDetailLocked[uuid]
	if !exist {
		return false
	}

	return isLocked
}

func updateNodeDetailLock(uuid string) {
	updateNodeDetailLocked[uuid] = true
}

func updateNodeDetailUnlock(uuid string) {
	updateNodeDetailLocked[uuid] = false
}

// makeRackNumber : Split IP address and add numbers of 4 sections with prefix length.
func makeRackNumber(bmcIPCIDR string) (int, error) {
	err := iputil.CheckCIDRStr(bmcIPCIDR)
	if err != nil {
		return 0, err
	}

	_, netIPNet, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		return 0, err
	}

	CIDRSplit := strings.Split(bmcIPCIDR, "/")
	netmask, err := strconv.Atoi(CIDRSplit[1])
	if err != nil {
		return 0, err
	}

	networkIPStr := netIPNet.IP.String()
	networkIPSplit := strings.Split(networkIPStr, ".")

	var ipSum = 0
	for _, split := range networkIPSplit {
		s, _ := strconv.Atoi(split)
		ipSum += s
	}

	rackNumber := ipSum + netmask

	return rackNumber, nil
}

func checkGroupIDExist(groupID int64) error {
	resGetGroupList, hccErrStack := client.RC.GetGroupList(&pb.Empty{})
	if hccErrStack != nil && (*hccErrStack.Stack())[0].Code() != 0 {
		return (*hccErrStack.Stack())[0].ToError()
	}

	for _, pGroup := range resGetGroupList.Group {
		if pGroup.Id == groupID {
			return nil
		}
	}

	return errors.New("given group ID is not in the database")
}

func checkNICSpeed(speed int) error {
	switch speed {
	case 10, 100, 1000, 2500, 5000, 10000, 20000, 40000:
		return nil
	default:
		return errors.New("unknown or not supported NIC speed")
	}
}

// DoUpdateAllNodes : Update the database of a specific node by getting bmcIP
func DoUpdateAllNodes(bmcIPCIDR string, wait *sync.WaitGroup, isNew bool, reqNode *pb.Node) (string, error) {
	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): Updating for bmc IP " + bmcIPCIDR)
	}

	if isNew {
		err := checkGroupIDExist(reqNode.GroupID)
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return "", err
		}

		err = checkNICSpeed(int(reqNode.NicSpeedMbps))
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return "", err
		}
	}

	rackNumber, err := makeRackNumber(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	netIP, _, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}
	bmcIP := netIP.String()

	serialNo, err := GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " Serial No: " + serialNo)
	}

	uuid, err := GetUUID(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " UUID: " + uuid)
	}

	bmcMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNumBMC), true)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " BMC MAC Addr: " + bmcMAC)
	}

	pxeMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNumPXE), false)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " PXE MAC Addr: " + pxeMAC)
	}

	processors, err := getNumCPU(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " Processors: " + strconv.Itoa(processors))
	}

	cpuCores, err := GetProcessorsCores(bmcIP, serialNo, processors)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " CPU Cores: " + strconv.Itoa(cpuCores))
	}

	memory, err := GetTotalSystemMemory(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " Memory: " + strconv.Itoa(memory))
	}

	node := pb.Node{
		UUID:       uuid,
		BmcMacAddr: bmcMAC,
		BmcIP:      bmcIPCIDR,
		PXEMacAddr: pxeMAC,
		CPUCores:   int32(cpuCores),
		Memory:     int32(memory),
		RackNumber: int32(rackNumber),
	}

	if isNew {
		sql := "insert into node(uuid, node_name, group_id, server_uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, " +
			"nic_speed_mbps, " +
			"description, rack_number, charge_cpu, charge_memory, charge_nic, created_at, available) " +
			"values (?, ?, ?, '', ?, ?, ?, '', ?, ?, " +
			"?, " +
			"?, ?, ?, ?, ?, now(), 1)"

		var stmt *dbsql.Stmt
		stmt, err := mysql.Prepare(sql)
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return "", err
		}
		defer func() {
			_ = stmt.Close()
		}()
		_, err = stmt.Exec(node.UUID, reqNode.NodeName, reqNode.GroupID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.CPUCores, node.Memory,
			reqNode.NicSpeedMbps,
			reqNode.GetDescription(), node.RackNumber, reqNode.ChargeCPU, reqNode.ChargeMemory, reqNode.ChargeNIC)
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return "", err
		}

		wait.Done()
		return uuid, nil
	}

	sql := "update node set uuid = ?, bmc_mac_addr = ?, pxe_mac_addr = ?, cpu_cores = ?, memory = ?, rack_number = ? where bmc_ip = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return "", err
	}

	result, err2 := stmt.Exec(node.UUID, node.BmcMacAddr, node.PXEMacAddr, node.CPUCores, node.Memory, node.RackNumber, node.BmcIP)
	if err2 != nil {
		logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err2.Error())
		_ = stmt.Close()
		wait.Done()
		return "", err2
	}
	_ = stmt.Close()

	if config.Ipmi.Debug == "on" {
		result, err := result.LastInsertId()
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
		} else {
			logger.Logger.Print("DoUpdateAllNodes(): " + bmcIPCIDR + " result=" + strconv.Itoa(int(result)))
		}
	}

	wait.Done()
	return uuid, nil
}

// UpdateNodesAll : Get all infos from IPMI nodes and update the database (except power state)
func UpdateNodesAll() {
	var bmcIPCIDR string

	sql := "select bmc_ip from node where available = 1"
	stmt, err := mysql.Query(sql)
	if err != nil {
		logger.Logger.Println("UpdateNodesAll(): err=" + err.Error())
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	resReadNodeNum, errCode, errText := daoext.ReadNodeNum(&pb.ReqGetNodeNum{})
	if errCode != 0 {
		logger.Logger.Println("UpdateNodesAll(): err=" + errText)
		return
	}

	var wait sync.WaitGroup
	wait.Add(int(resReadNodeNum.Num))

	for stmt.Next() {
		err := stmt.Scan(&bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateNodesAll(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			continue
		}

		go func(bmcIPCIDR string, wait *sync.WaitGroup) {
			_, _ = DoUpdateAllNodes(bmcIPCIDR, wait, false, nil)
		}(bmcIPCIDR, &wait)
	}

	wait.Wait()

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("UpdateNodesAll(): Done")
	}
}

// DoUpdateStatusNodes : Get status from a specific node and update the database
func DoUpdateStatusNodes(uuid interface{}, bmcIPCIDR string, oldStatus string, wait *sync.WaitGroup) error {
	err := iputil.CheckCIDRStr(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	netIP, _, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}
	bmcIP := netIP.String()

	if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + "'s UUID is currently empty. Skipping...")
		}
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateStatusNodes(): Updating for bmc IP " + bmcIPCIDR)
	}

	serialNo, err := GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " Serial No: " + serialNo)
	}

	powerState, err := GetPowerState(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " Power State: " + powerState)
	}

	_, errCode, _ := daoext.ReadNodeDetail(fmt.Sprintf("%s", uuid))
	if errCode == hcc_errors.FluteSQLNoResult {
		ScheduleUpdateNodeDetail(fmt.Sprintf("%s", uuid))
	} else if oldStatus == powerState {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR +
				" No need to update power state (Power State: " + powerState + ")")
		}
		wait.Done()
		return nil
	} else {
		ScheduleUpdateNodeDetail(fmt.Sprintf("%s", uuid))
	}

	node := pb.Node{
		UUID:   fmt.Sprintf("%s", uuid),
		Status: powerState,
	}

	sql := "update node set status = ? where uuid = ?"
	stmt, err := mysql.Prepare(sql)
	if err != nil {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	result, err2 := stmt.Exec(node.Status, node.UUID)
	if err2 != nil {
		logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err2.Error())
		_ = stmt.Close()
		wait.Done()
		return err2
	}
	_ = stmt.Close()

	if config.Ipmi.Debug == "on" {
		result, err := result.LastInsertId()
		if err != nil {
			logger.Logger.Println("DoUpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
		} else {
			logger.Logger.Print("DoUpdateStatusNodes(): " + bmcIPCIDR + " result=" + strconv.Itoa(int(result)))
		}
	}

	wait.Done()
	return nil
}

// UpdateNodesStatus : Get status from IPMI nodes and update the database
func UpdateNodesStatus() {
	var uuid interface{}
	var bmcIPCIDR string
	var oldStatus string

	sql := "select uuid, bmc_ip, status from node where available = 1"
	stmt, err := mysql.Query(sql)
	if err != nil {
		logger.Logger.Println("UpdateNodesStatus(): err=" + err.Error())
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	resReadNodeNum, errCode, errText := daoext.ReadNodeNum(&pb.ReqGetNodeNum{})
	if errCode != 0 {
		logger.Logger.Println("UpdateNodesStatus(): err=" + errText)
		return
	}

	var wait sync.WaitGroup
	wait.Add(int(resReadNodeNum.Num))

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIPCIDR, &oldStatus)
		if err != nil {
			logger.Logger.Println("UpdateNodesStatus(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			continue
		}

		go func(uuid interface{}, bmcIPCIDR string, wait *sync.WaitGroup) {
			_ = DoUpdateStatusNodes(uuid, bmcIPCIDR, oldStatus, wait)
		}(uuid, bmcIPCIDR, &wait)
	}

	wait.Wait()

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("UpdateNodesStatus(): Done")
	}
}

// UpdateServerStatus : Update status of the server
func UpdateServerStatus() {
	var status string

	resGetServerList, err := client.RC.GetServerList(&pb.ReqGetServerList{})
	if err != nil {
		logger.Logger.Println("UpdateServerStatus(): err=" + err.Error())
		return
	}

	for _, server := range resGetServerList.Server {
		if strings.ToLower(server.Status) == "creating" {
			continue
		}

		sql := "select status from node where server_uuid = '" + server.UUID + "'"
		stmt, err := mysql.Query(sql)
		if err != nil {
			logger.Logger.Println("UpdateServerStatus(): err=" + err.Error())
			return
		}
		defer func() {
			_ = stmt.Close()
		}()

		var isAllTurnedOff = true

		for stmt.Next() {
			err := stmt.Scan(&status)
			if err != nil {
				logger.Logger.Println("UpdateServerStatus(): err=" + err.Error())
				continue
			}

			if strings.ToLower(status) != "off" {
				isAllTurnedOff = false
				break
			}
		}

		if isAllTurnedOff {
			_, err := client.RC.UpdateServer(&pb.ReqUpdateServer{
				Server: &pb.Server{
					UUID:   server.UUID,
					Status: "Stopped",
				},
			})
			if err != nil {
				logger.Logger.Println("UpdateServerStatus(): err=" + err.Error())
			}
		}
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("UpdateServerStatus(): Done")
	}
}

func doUpdateNodesDetail(uuid interface{}, bmcIPCIDR string, wait *sync.WaitGroup) error {
	err := iputil.CheckCIDRStr(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	netIP, _, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}
	bmcIP := netIP.String()

	if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + "'s UUID is currently empty. Skipping...")
		}
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("doUpdateNodesDetail(): Updating for bmc IP " + bmcIPCIDR)
	}

	serialNo, err := GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " Serial No: " + serialNo)
	}

	nodeDetailData, err := GetNodeDetailData(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("doUpdateNodesDetail(): "+bmcIPCIDR+" NodeDetailData : ", nodeDetailData)
	}

	nodeUUID := fmt.Sprintf("%s", uuid)
	nodeDetail := pb.NodeDetail{
		NodeUUID:       nodeUUID,
		NodeDetailData: nodeDetailData,
	}

	sql := "select node_uuid from node_detail where node_uuid = ?"
	row := mysql.Db.QueryRow(sql, uuid)
	err = mysql.QueryRowScan(row, &uuid)
	if err != nil {
		logger.Logger.Println("doUpdateNodesDetail(): Inserting not existing new node_detail")

		sql = "insert into node_detail(node_uuid, node_detail_data) values (?, ?)"
		stmt, err := mysql.Prepare(sql)
		if err != nil {
			logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return err
		}

		result, err2 := stmt.Exec(nodeDetail.NodeUUID, nodeDetail.NodeDetailData)
		if err2 != nil {
			logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err2.Error())
			_ = stmt.Close()
			wait.Done()
			return err2
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			} else {
				logger.Logger.Print("doUpdateNodesDetail(): " + bmcIPCIDR + " result=" + strconv.Itoa(int(result)))
			}
		}
	} else {
		sql = "update node_detail set node_detail_data = ? where node_uuid = ?"
		stmt, err := mysql.Prepare(sql)
		if err != nil {
			logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return err
		}

		result, err2 := stmt.Exec(nodeDetail.NodeDetailData, nodeDetail.NodeUUID)
		if err2 != nil {
			logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err2.Error())
			_ = stmt.Close()
			wait.Done()
			return err2
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Println("doUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			} else {
				logger.Logger.Print("doUpdateNodesDetail(): " + bmcIPCIDR + " result=" + strconv.Itoa(int(result)))
			}
		}
	}

	wait.Done()
	return nil
}

func updateNodeDetail(uuid string) error {
	var bmcIPCIDR string

	sql := "select uuid, bmc_ip from node where available = 1 and uuid = '" + uuid + "'"
	stmt, err := mysql.Query(sql)
	if err != nil {
		logger.Logger.Println("updateNodeDetail(): uuid=" + uuid + " err=" + err.Error())
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var wait sync.WaitGroup
	wait.Add(1)

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("updateNodeDetail(): " + bmcIPCIDR + " uuid=" + uuid + " err=" + err.Error())
			continue
		}

		err = doUpdateNodesDetail(uuid, bmcIPCIDR, &wait)
		if err != nil {
			logger.Logger.Println("updateNodeDetail(): uuid=" + uuid + " err=" + err.Error())
			return err
		}
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("updateNodeDetail(): Done of uuid=" + uuid)
	}

	return nil
}

func queueCheckNodeAll() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckNodeAll(): Queued of running CheckNodeAll() after " + strconv.Itoa(int(config.Ipmi.CheckNodeAllIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckNodeAllIntervalMs))
		CheckNodeAll()
	}()
}

func queueCheckNodeStatus() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckNodeStatus(): Queued of running CheckNodeStatus() after " + strconv.Itoa(int(config.Ipmi.CheckNodeStatusIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckNodeStatusIntervalMs))
		CheckNodeStatus()
	}()
}

func queueCheckServerStatus() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckServerStatus(): Queued of running CheckServerStatus() after " + strconv.Itoa(int(config.Ipmi.CheckServerStatusIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckServerStatusIntervalMs))
		CheckServerStatus()
	}()
}

// CheckNodeAll : Check node's all infos by interval of 'check_node_all_interval_ms' config option
func CheckNodeAll() {
	if checkNodeAllLocked {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckNodeAll(): Locked")
		}
		for true {
			if !checkNodeAllLocked {
				break
			}
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("CheckNodeAll(): Rerun after " + strconv.Itoa(int(config.Ipmi.CheckNodeAllIntervalMs)) + "ms")
			}
			delayMillisecond(time.Duration(config.Ipmi.CheckNodeAllIntervalMs))
		}
	}

	go func() {
		checkNodeAllLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckNodeAll(): Running UpdateNodesAll()")
		}
		UpdateNodesAll()
		checkNodeAllUnlock()
	}()

	queueCheckNodeAll()
}

// CheckNodeStatus : Check node's power status by interval of 'check_node_status_interval_ms' config option
func CheckNodeStatus() {
	if checkNodeStatusLocked {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckNodeStatus(): Locked")
		}
		for true {
			if !checkNodeStatusLocked {
				break
			}
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("CheckNodeStatus(): Rerun after " + strconv.Itoa(int(config.Ipmi.CheckNodeStatusIntervalMs)) + "ms")
			}
			delayMillisecond(time.Duration(config.Ipmi.CheckNodeStatusIntervalMs))
		}
	}

	go func() {
		checkNodeStatusLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckNodeStatus(): Running UpdateNodesStatus()")
		}
		UpdateNodesStatus()
		checkNodeStatusUnlock()
	}()

	queueCheckNodeStatus()
}

// CheckServerStatus : Check server's power status by interval of 'check_server_status_interval_ms' config option
func CheckServerStatus() {
	if checkServerStatusLocked {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckServerStatus(): Locked")
		}
		for true {
			if !checkServerStatusLocked {
				break
			}
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("CheckServerStatus(): Rerun after " + strconv.Itoa(int(config.Ipmi.CheckServerStatusIntervalMs)) + "ms")
			}
			delayMillisecond(time.Duration(config.Ipmi.CheckServerStatusIntervalMs))
		}
	}

	go func() {
		checkServerStatusLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckServerStatus(): Running UpdateServerStatus()")
		}
		UpdateServerStatus()
		checkServerStatusUnlock()
	}()

	queueCheckServerStatus()
}

// ScheduleUpdateNodeDetail : Schedule of updating detail infos of the node
func ScheduleUpdateNodeDetail(uuid string) {
	if isUpdateNodeDetailLocked(uuid) {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("ScheduleUpdateNodeDetail(): Locked for uuid=" + uuid)
		}
		for true {
			if !isUpdateNodeDetailLocked(uuid) {
				break
			}
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("ScheduleUpdateNodeDetail(): Rerun after " +
					strconv.Itoa(int(config.Ipmi.UpdateNodeDetailRetryIntervalMs)) + "ms for uuid=" + uuid)
			}
			delayMillisecond(time.Duration(config.Ipmi.UpdateNodeDetailRetryIntervalMs))
		}
	}

	go func() {
		updateNodeDetailLock(uuid)
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("ScheduleUpdateNodeDetail(): Running updateNodeDetail() for uuid=" + uuid)
		}
		err := updateNodeDetail(uuid)
		updateNodeDetailUnlock(uuid)
		if err != nil {
			ScheduleUpdateNodeDetail(uuid)
		}
	}()
}
