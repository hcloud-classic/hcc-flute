package ipmi

import (
	"hcc/flute/config"
	"hcc/flute/logger"
	"hcc/flute/mysql"
	"hcc/flute/types"
	"strconv"
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

// UpdateAllNodes : Get all infos from IPMI nodes and update database (except power state)
func UpdateAllNodes() (interface{}, error) {
	var nodes []types.Node
	var bmcIP string

	sql := "select bmc_ip from node where active = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		uuid, err := GetUUID(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		bmcMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoBMC), true)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		pxeMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoPXE), false)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		cpuCores, err := GetProcessorsCores(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		memory, err := GetTotalSystemMemory(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		node := types.Node{
			UUID:       uuid,
			BmcMacAddr: bmcMAC,
			BmcIP:      bmcIP,
			PXEMacAddr: pxeMAC,
			CPUCores:   cpuCores,
			Memory:     memory,
		}

		sql := "update node set bmc_mac_addr = ?, pxe_mac_addr = ?, cpu_cores = ?, memory = ? where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, nil
		}

		result, err2 := stmt.Exec(node.BmcMacAddr, node.PXEMacAddr, node.CPUCores, node.Memory, node.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, nil
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println(result.LastInsertId())
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateStatusNodes : Get status from IPMI nodes and update database
func UpdateStatusNodes() (interface{}, error) {
	var nodes []types.Node
	var uuid string
	var bmcIP string

	sql := "select uuid, bmc_ip from node where active = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		powerState, err := GetPowerState(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		node := types.Node{
			UUID:   uuid,
			Status: powerState,
		}

		sql = "update node set status = ? where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, nil
		}

		result, err2 := stmt.Exec(node.Status, node.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, nil
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println(result.LastInsertId())
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateNodesDetail : Get detail infos from IPMI nodes and update database
func UpdateNodesDetail() (interface{}, error) {
	var nodedetails []types.NodeDetail
	var uuid string
	var bmcIP string

	sql := "select uuid, bmc_ip from node where active = 1"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		model, err := GetProcessorModel(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		threads, err := GetProcessorsThreads(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		nodedetail := types.NodeDetail{
			NodeUUID:      uuid,
			CPUModel:      model,
			CPUProcessors: processors,
			CPUThreads:    threads,
		}

		sql := "select node_uuid from node_detail where node_uuid = ?"
		err = mysql.Db.QueryRow(sql, uuid).Scan(&uuid)
		logger.Logger.Println(err)
		if err != nil {
			logger.Logger.Println("Inserting not existing new node_detail")

			sql = "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
			stmt, err := mysql.Db.Prepare(sql)
			if err != nil {
				logger.Logger.Println(err.Error())
				return nil, nil
			}
			defer func() {
				_ = stmt.Close()
			}()
			result, err2 := stmt.Exec(nodedetail.NodeUUID, nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads)
			if err2 != nil {
				logger.Logger.Println(err2)
				return nil, nil
			}
			logger.Logger.Println(result.LastInsertId())
		} else {
			sql = "update node_detail set cpu_model = ?, cpu_processors = ?, cpu_threads = ? where node_uuid = ?"
			stmt, err := mysql.Db.Prepare(sql)
			if err != nil {
				logger.Logger.Println(err.Error())
				return nil, nil
			}

			result, err2 := stmt.Exec(nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads, nodedetail.NodeUUID)
			if err2 != nil {
				logger.Logger.Println(err2)
				return nil, nil
			}
			_ = stmt.Close()

			if config.Ipmi.Debug == "on" {
				logger.Logger.Println(result.LastInsertId())
			}
		}
		defer func() {
			_ = stmt.Close()
		}()

		nodedetails = append(nodedetails, nodedetail)
	}

	return nodedetails, nil
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
		_, _ = UpdateAllNodes()
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
		_, _ = UpdateStatusNodes()
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
		_, _ = UpdateNodesDetail()
		checkNodesDetailUnlock()
	}()

	queueNodesDetail()
}
