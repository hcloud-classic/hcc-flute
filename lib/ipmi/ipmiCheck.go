package ipmi

import (
	"fmt"
	pb "hcc/flute/action/grpc/pb/rpcflute"
	"hcc/flute/daoext"
	"hcc/flute/lib/config"
	"hcc/flute/lib/iputil"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var checkAllLocked = false
var checkStatusLocked = false
var checkNodesDetailLocked = false

func delayMillisecond(n time.Duration) {
	time.Sleep(n * time.Millisecond)
}

func checkAllLock() {
	checkAllLocked = true
}

func checkAllUnlock() {
	checkAllLocked = false
}

func checkStatusLock() {
	checkStatusLocked = true
}

func checkStatusUnlock() {
	checkStatusLocked = false
}

func checkNodesDetailLock() {
	checkNodesDetailLocked = true
}

func checkNodesDetailUnlock() {
	checkNodesDetailLocked = false
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

// DoUpdateAllNodes : Update the database of a specific node by getting bmcIP
func DoUpdateAllNodes(bmcIPCIDR string, wait *sync.WaitGroup, isNew bool, description string) (string, error) {
	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateAllNodes(): Updating for bmc IP " + bmcIPCIDR)
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

	processors, err := GetProcessors(bmcIP, serialNo)
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
		sql := "insert into node(uuid, server_uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, rack_number, created_at, available) values (?, '', ?, ?, ?, '', ?, ?, ?, ?, now(), 1)"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return "", err
		}
		defer func() {
			_ = stmt.Close()
		}()
		_, err = stmt.Exec(node.UUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.CPUCores, node.Memory, description, node.RackNumber)
		if err != nil {
			logger.Logger.Println("DoUpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return "", err
		}

		wait.Done()
		return uuid, nil
	}

	sql := "update node set uuid = ?, bmc_mac_addr = ?, pxe_mac_addr = ?, cpu_cores = ?, memory = ?, rack_number = ? where bmc_ip = ?"
	stmt, err := mysql.Db.Prepare(sql)
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

// UpdateAllNodes : Get all infos from IPMI nodes and update the database (except power state)
func UpdateAllNodes() {
	var bmcIPCIDR string

	sql := "select bmc_ip from node where available = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println("UpdateAllNodes(): err=" + err.Error())
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	resReadNodeNum, errCode, errText := daoext.ReadNodeNum()
	if errCode != 0 {
		logger.Logger.Println("UpdateAllNodes(): err=" + errText)
		return
	}

	var wait sync.WaitGroup
	wait.Add(int(resReadNodeNum.Num))

	for stmt.Next() {
		err := stmt.Scan(&bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIPCIDR + " err=" + err.Error())
			continue
		}

		go func(bmcIPCIDR string, wait *sync.WaitGroup) {
			_, _ = DoUpdateAllNodes(bmcIPCIDR, wait, false, "")
		}(bmcIPCIDR, &wait)
	}

	wait.Wait()

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("UpdateAllNodes(): Done")
	}
}

// DoUpdateStatusNodes : Get status from a specific node and update the database
func DoUpdateStatusNodes(uuid interface{}, bmcIPCIDR string, wait *sync.WaitGroup) error {
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

	node := pb.Node{
		UUID:   fmt.Sprintf("%s", uuid),
		Status: powerState,
	}

	sql := "update node set status = ? where uuid = ?"
	stmt, err := mysql.Db.Prepare(sql)
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

// UpdateStatusNodes : Get status from IPMI nodes and update the database
func UpdateStatusNodes() {
	var uuid interface{}
	var bmcIPCIDR string

	sql := "select uuid, bmc_ip from node where available = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println("UpdateStatusNodes(): err=" + err.Error())
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	resReadNodeNum, errCode, errText := daoext.ReadNodeNum()
	if errCode != 0 {
		logger.Logger.Println("UpdateStatusNodes(): err=" + errText)
		return
	}

	var wait sync.WaitGroup
	wait.Add(int(resReadNodeNum.Num))

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateStatusNodes(): " + bmcIPCIDR + " err=" + err.Error())
			continue
		}

		go func(uuid interface{}, bmcIPCIDR string, wait *sync.WaitGroup) {
			_ = DoUpdateStatusNodes(uuid, bmcIPCIDR, wait)
		}(uuid, bmcIPCIDR, &wait)
	}

	wait.Wait()

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("UpdateStatusNodes(): Done")
	}
}

// DoUpdateNodesDetail : Get detail infos from a specific node and update the database
func DoUpdateNodesDetail(uuid interface{}, bmcIPCIDR string, wait *sync.WaitGroup) error {
	err := iputil.CheckCIDRStr(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	netIP, _, err := net.ParseCIDR(bmcIPCIDR)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}
	bmcIP := netIP.String()

	if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + "'s UUID is currently empty. Skipping...")
		}
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateNodesDetail(): Updating for bmc IP " + bmcIPCIDR)
	}

	serialNo, err := GetSerialNo(bmcIP)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " Serial No: " + serialNo)
	}

	processorModel, err := GetProcessorModel(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " Processor Model: " + processorModel)
	}

	processors, err := GetProcessors(bmcIP, serialNo)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " Processors : " + strconv.Itoa(processors))
	}

	threads, err := GetProcessorsThreads(bmcIP, serialNo, processors)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
		wait.Done()
		return err
	}

	if config.Ipmi.Debug == "on" {
		logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " Threads : " + strconv.Itoa(threads))
	}

	nodeUUID := fmt.Sprintf("%s", uuid)
	nodeDetail := pb.NodeDetail{
		NodeUUID:      nodeUUID,
		CPUModel:      processorModel,
		CPUProcessors: int32(processors),
		CPUThreads:    int32(threads),
	}

	sql := "select node_uuid from node_detail where node_uuid = ?"
	err = mysql.Db.QueryRow(sql, uuid).Scan(&uuid)
	if err != nil {
		logger.Logger.Println("DoUpdateNodesDetail(): Inserting not existing new node_detail")

		sql = "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return err
		}

		result, err2 := stmt.Exec(nodeDetail.NodeUUID, nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads)
		if err2 != nil {
			logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err2.Error())
			_ = stmt.Close()
			wait.Done()
			return err2
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			} else {
				logger.Logger.Print("DoUpdateNodesDetail(): " + bmcIPCIDR + " result=" + strconv.Itoa(int(result)))
			}
		}
	} else {
		sql = "update node_detail set cpu_model = ?, cpu_processors = ?, cpu_threads = ? where node_uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			wait.Done()
			return err
		}

		result, err2 := stmt.Exec(nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads, nodeDetail.NodeUUID)
		if err2 != nil {
			logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err2.Error())
			_ = stmt.Close()
			wait.Done()
			return err2
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Println("DoUpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			} else {
				logger.Logger.Print("DoUpdateNodesDetail(): " + bmcIPCIDR + " result=" + strconv.Itoa(int(result)))
			}
		}
	}

	wait.Done()
	return nil
}

// UpdateNodesDetail : Get detail infos from IPMI nodes and update the database
func UpdateNodesDetail() {
	var uuid interface{}
	var bmcIPCIDR string

	sql := "select uuid, bmc_ip from node where available = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println("UpdateNodesDetail(): err=" + err.Error())
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	resReadNodeNum, errCode, errText := daoext.ReadNodeNum()
	if errCode != 0 {
		logger.Logger.Println("UpdateNodesDetail(): err=" + errText)
		return
	}

	var wait sync.WaitGroup
	wait.Add(int(resReadNodeNum.Num))

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIPCIDR + " err=" + err.Error())
			continue
		}

		go func(uuid interface{}, bmcIPCIDR string, wait *sync.WaitGroup) {
			_ = DoUpdateNodesDetail(uuid, bmcIPCIDR, wait)
		}(uuid, bmcIPCIDR, &wait)
	}

	wait.Wait()

	if config.Ipmi.Debug == "on" {
		logger.Logger.Print("UpdateNodesDetail(): Done")
	}
}

func queueCheckAll() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckAll(): Rerun CheckAll() after " + strconv.Itoa(int(config.Ipmi.CheckAllIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckAllIntervalMs))
		CheckAll()
	}()
}

func queueCheckStatus() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckStatus(): Rerun CheckStatus() after " + strconv.Itoa(int(config.Ipmi.CheckStatusIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckStatusIntervalMs))
		CheckStatus()
	}()
}

func queueNodesDetail() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueNodesDetail(): Rerun NodesDetail() after " + strconv.Itoa(int(config.Ipmi.CheckNodesDetailIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckNodesDetailIntervalMs))
		CheckNodesDetail()
	}()
}

// CheckAll : Check all IPMI infos by 'check_all_interval_ms' config option
func CheckAll() {
	if checkAllLocked {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckAll(): Locked")
		}
		queueCheckAll()
		return
	}

	go func() {
		checkAllLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckAll(): Running UpdateAllNodes()")
		}
		UpdateAllNodes()
		checkAllUnlock()
	}()

	queueCheckAll()
}

// CheckStatus : Check power status of IPMI nodes by 'check_status_interval_ms' config option
func CheckStatus() {
	if checkStatusLocked {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckStatus(): Locked")
		}
		queueCheckStatus()
		return
	}

	go func() {
		checkStatusLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckStatus(): Running UpdateStatusNodes()")
		}
		UpdateStatusNodes()
		checkStatusUnlock()
	}()

	queueCheckStatus()
}

// CheckNodesDetail : Check detail infos of IPMI nodes by 'check_nodes_detail_interval_ms' config option
func CheckNodesDetail() {
	if checkNodesDetailLocked {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckNodesDetail(): Locked")
		}
		queueNodesDetail()
		return
	}

	go func() {
		checkNodesDetailLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("CheckNodesDetail(): Running UpdateNodesDetail()")
		}
		UpdateNodesDetail()
		checkNodesDetailUnlock()
	}()

	queueNodesDetail()
}
