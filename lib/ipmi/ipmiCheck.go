package ipmi

import (
	"fmt"
	pb "hcc/flute/action/grpc/pb/rpcflute"
	"hcc/flute/lib/config"
	"hcc/flute/lib/iputil"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"net"
	"strconv"
	"strings"
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
	fmt.Println(networkIPStr)
	networkIPSplit := strings.Split(networkIPStr, ".")

	var ipSum = 0
	for _, split := range networkIPSplit {
		s, _ := strconv.Atoi(split)
		ipSum += s
	}

	rackNumber := ipSum + netmask

	return rackNumber, nil
}

// UpdateAllNodes : Get all infos from IPMI nodes and update database (except power state)
func UpdateAllNodes() {
	var bmcIPCIDR string

	sql := "select bmc_ip from node where available = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&bmcIPCIDR)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): Updating for bmc IP " + bmcIPCIDR)
		}

		rackNumber, err := makeRackNumber(bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateAllNodes(): " + err.Error())
			continue
		}

		netIP, _, err := net.ParseCIDR(bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateAllNodes(): " + err.Error())
			continue
		}
		bmcIP := netIP.String()

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " Serial No: " + serialNo)
		}

		uuid, err := GetUUID(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " UUID: " + uuid)
		}

		bmcMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNumBMC), true)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " BMC MAC Addr: " + bmcMAC)
		}

		pxeMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNumPXE), false)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " PXE MAC Addr: " + pxeMAC)
		}

		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " Processors: " + strconv.Itoa(processors))
		}

		cpuCores, err := GetProcessorsCores(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " CPU Cores: " + strconv.Itoa(cpuCores))
		}

		memory, err := GetTotalSystemMemory(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " Memory: " + strconv.Itoa(memory))
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

		sql := "update node set uuid = ?, bmc_mac_addr = ?, pxe_mac_addr = ?, cpu_cores = ?, memory = ?, rack_number = ? where bmc_ip = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		result, err2 := stmt.Exec(node.UUID, node.BmcMacAddr, node.PXEMacAddr, node.CPUCores, node.Memory, node.RackNumber, node.BmcIP)
		if err2 != nil {
			logger.Logger.Println(err2)
			_ = stmt.Close()
			continue
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Print("UpdateAllNodes(): err=" + err.Error())
			} else {
				logger.Logger.Print("UpdateAllNodes(): result=" + strconv.Itoa(int(result)))
			}
		}
	}
}

// UpdateStatusNodes : Get status from IPMI nodes and update database
func UpdateStatusNodes() {
	var uuid interface{}
	var bmcIPCIDR string

	sql := "select uuid, bmc_ip from node where available = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIPCIDR)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		err = iputil.CheckCIDRStr(bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateStatusNodes(): " + err.Error())
			return
		}

		netIP, _, err := net.ParseCIDR(bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateStatusNodes(): " + err.Error())
			return
		}
		bmcIP := netIP.String()

		if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("UpdateAllNodes(): " + bmcIP + "'s UUID is currently empty. Skipping...")
			}
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateStatusNodes(): Updating for bmc IP " + bmcIP)
		}

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateStatusNodes(): " + bmcIP + " Serial No: " + serialNo)
		}

		powerState, err := GetPowerState(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateStatusNodes(): " + bmcIP + " Power State: " + powerState)
		}

		node := pb.Node{
			UUID:   fmt.Sprintf("%s", uuid),
			Status: powerState,
		}

		sql = "update node set status = ? where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		result, err2 := stmt.Exec(node.Status, node.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			_ = stmt.Close()
			continue
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Print("UpdateStatusNodes(): err=" + err.Error())
			} else {
				logger.Logger.Print("UpdateStatusNodes(): result=" + strconv.Itoa(int(result)))
			}
		}
	}
}

// UpdateNodesDetail : Get detail infos from IPMI nodes and update database
func UpdateNodesDetail() {
	var uuid interface{}
	var bmcIPCIDR string

	sql := "select uuid, bmc_ip from node where available = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIPCIDR)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		err = iputil.CheckCIDRStr(bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateNodesDetail(): " + err.Error())
			return
		}

		netIP, _, err := net.ParseCIDR(bmcIPCIDR)
		if err != nil {
			logger.Logger.Println("UpdateNodesDetail(): " + err.Error())
			return
		}
		bmcIP := netIP.String()

		if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + "'s UUID is currently empty. Skipping...")
			}
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): Updating for bmc IP " + bmcIP)
		}

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Serial No: " + serialNo)
		}

		processorModel, err := GetProcessorModel(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Processor Model: " + processorModel)
		}

		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Processors : " + strconv.Itoa(processors))
		}

		threads, err := GetProcessorsThreads(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Threads : " + strconv.Itoa(threads))
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
			logger.Logger.Println("UpdateNodesDetail(): Inserting not existing new node_detail")

			sql = "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
			stmt, err := mysql.Db.Prepare(sql)
			if err != nil {
				logger.Logger.Println(err)
				continue
			}

			result, err2 := stmt.Exec(nodeDetail.NodeUUID, nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads)
			if err2 != nil {
				logger.Logger.Println(err2)
				_ = stmt.Close()
				continue
			}
			_ = stmt.Close()

			if config.Ipmi.Debug == "on" {
				result, err := result.LastInsertId()
				if err != nil {
					logger.Logger.Print("UpdateNodesDetail(): err=" + err.Error())
				} else {
					logger.Logger.Print("UpdateNodesDetail(): result=" + strconv.Itoa(int(result)))
				}
			}
		} else {
			sql = "update node_detail set cpu_model = ?, cpu_processors = ?, cpu_threads = ? where node_uuid = ?"
			stmt, err := mysql.Db.Prepare(sql)
			if err != nil {
				logger.Logger.Println(err)
				continue
			}

			result, err2 := stmt.Exec(nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads, nodeDetail.NodeUUID)
			if err2 != nil {
				logger.Logger.Println(err2)
				_ = stmt.Close()
				continue
			}
			_ = stmt.Close()

			if config.Ipmi.Debug == "on" {
				result, err := result.LastInsertId()
				if err != nil {
					logger.Logger.Print("UpdateNodesDetail(): err=" + err.Error())
				} else {
					logger.Logger.Print("UpdateNodesDetail(): result=" + strconv.Itoa(int(result)))
				}
			}
		}
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
			logger.Logger.Println("NodesDetail(): Locked")
		}
		queueNodesDetail()
		return
	}

	go func() {
		checkNodesDetailLock()
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("NodesDetail(): Running UpdateNodesDetail()")
		}
		UpdateNodesDetail()
		checkNodesDetailUnlock()
	}()

	queueNodesDetail()
}
