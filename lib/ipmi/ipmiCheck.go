package ipmi

import (
<<<<<<< HEAD
	"fmt"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
=======
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
>>>>>>> f41ff24 (Refactoring packages structure)
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
<<<<<<< HEAD
	var nodes []model.Node
=======
	var nodes []types.Node
>>>>>>> f41ff24 (Refactoring packages structure)
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
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): Updating for bmc IP " + bmcIP)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " Serial No: " + serialNo)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		uuid, err := GetUUID(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " UUID: " + uuid)
		}

		bmcMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNumBMC), true)
=======
		bmcMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoBMC), true)
>>>>>>> f41ff24 (Refactoring packages structure)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " BMC MAC Addr: " + bmcMAC)
		}

		pxeMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNumPXE), false)
=======
		pxeMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoPXE), false)
>>>>>>> f41ff24 (Refactoring packages structure)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " PXE MAC Addr: " + pxeMAC)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " Processors: " + strconv.Itoa(processors))
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		cpuCores, err := GetProcessorsCores(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " CPU Cores: " + strconv.Itoa(cpuCores))
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		memory, err := GetTotalSystemMemory(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateAllNodes(): " + bmcIP + " Memory: " + strconv.Itoa(memory))
		}

		node := model.Node{
=======
		node := types.Node{
>>>>>>> f41ff24 (Refactoring packages structure)
			UUID:       uuid,
			BmcMacAddr: bmcMAC,
			BmcIP:      bmcIP,
			PXEMacAddr: pxeMAC,
			CPUCores:   cpuCores,
			Memory:     memory,
		}

<<<<<<< HEAD
		sql := "update node set uuid = ?, bmc_mac_addr = ?, pxe_mac_addr = ?, cpu_cores = ?, memory = ? where bmc_ip = ?"
=======
		sql := "update node set bmc_mac_addr = ?, pxe_mac_addr = ?, cpu_cores = ?, memory = ? where uuid = ?"
>>>>>>> f41ff24 (Refactoring packages structure)
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		result, err2 := stmt.Exec(node.UUID, node.BmcMacAddr, node.PXEMacAddr, node.CPUCores, node.Memory, node.BmcIP)
=======
		result, err2 := stmt.Exec(node.BmcMacAddr, node.PXEMacAddr, node.CPUCores, node.Memory, node.UUID)
>>>>>>> f41ff24 (Refactoring packages structure)
		if err2 != nil {
			logger.Logger.Println(err2)
			_ = stmt.Close()
			continue
		}
		_ = stmt.Close()

		if config.Ipmi.Debug == "on" {
<<<<<<< HEAD
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Print("UpdateAllNodes(): err=" + err.Error())
			} else {
				logger.Logger.Print("UpdateAllNodes(): result=" + strconv.Itoa(int(result)))
			}
=======
			logger.Logger.Println(result.LastInsertId())
>>>>>>> f41ff24 (Refactoring packages structure)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateStatusNodes : Get status from IPMI nodes and update database
func UpdateStatusNodes() (interface{}, error) {
<<<<<<< HEAD
	var nodes []model.Node
	var uuid interface{}
=======
	var nodes []types.Node
	var uuid string
>>>>>>> f41ff24 (Refactoring packages structure)
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
			continue
		}

<<<<<<< HEAD
		if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("UpdateAllNodes(): " + bmcIP + "'s UUID is currently empty. Skipping...")
			}
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateStatusNodes(): Updating for bmc IP " + bmcIP)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateStatusNodes(): " + bmcIP + " Serial No: " + serialNo)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		powerState, err := GetPowerState(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateStatusNodes(): " + bmcIP + " Power State: " + powerState)
		}

		node := model.Node {
			UUID:   fmt.Sprintf("%s", uuid),
=======
		node := types.Node{
			UUID:   uuid,
>>>>>>> f41ff24 (Refactoring packages structure)
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
<<<<<<< HEAD
			result, err := result.LastInsertId()
			if err != nil {
				logger.Logger.Print("UpdateStatusNodes(): err=" + err.Error())
			} else {
				logger.Logger.Print("UpdateStatusNodes(): result=" + strconv.Itoa(int(result)))
			}
=======
			logger.Logger.Println(result.LastInsertId())
>>>>>>> f41ff24 (Refactoring packages structure)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateNodesDetail : Get detail infos from IPMI nodes and update database
func UpdateNodesDetail() (interface{}, error) {
<<<<<<< HEAD
	var nodedetails []model.NodeDetail
	var uuid interface{}
=======
	var nodedetails []types.NodeDetail
	var uuid string
>>>>>>> f41ff24 (Refactoring packages structure)
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
			continue
		}

<<<<<<< HEAD
		if uuid == nil || len(fmt.Sprintf("%s", uuid)) == 0 {
			if config.Ipmi.Debug == "on" {
				logger.Logger.Println("UpdateAllNodes(): " + bmcIP + "'s UUID is currently empty. Skipping...")
			}
			continue
		}

		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): Updating for bmc IP " + bmcIP)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Serial No: " + serialNo)
		}

		processorModel, err := GetProcessorModel(bmcIP, serialNo)
=======
		model, err := GetProcessorModel(bmcIP, serialNo)
>>>>>>> f41ff24 (Refactoring packages structure)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Processor Model: " + processorModel)
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Processors : " + strconv.Itoa(processors))
		}

=======
>>>>>>> f41ff24 (Refactoring packages structure)
		threads, err := GetProcessorsThreads(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

<<<<<<< HEAD
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("UpdateNodesDetail(): " + bmcIP + " Threads : " + strconv.Itoa(threads))
		}

		nodeUUID := fmt.Sprintf("%s", uuid)
		nodeDetail := model.NodeDetail{
			NodeUUID:      nodeUUID,
			CPUModel:      processorModel,
=======
		nodedetail := types.NodeDetail{
			NodeUUID:      uuid,
			CPUModel:      model,
>>>>>>> f41ff24 (Refactoring packages structure)
			CPUProcessors: processors,
			CPUThreads:    threads,
		}

		sql := "select node_uuid from node_detail where node_uuid = ?"
		err = mysql.Db.QueryRow(sql, uuid).Scan(&uuid)
		if err != nil {
<<<<<<< HEAD
			logger.Logger.Println("UpdateNodesDetail(): Inserting not existing new node_detail")
=======
			logger.Logger.Println("Inserting not existing new node_detail")
>>>>>>> f41ff24 (Refactoring packages structure)

			sql = "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
			stmt, err := mysql.Db.Prepare(sql)
			if err != nil {
				logger.Logger.Println(err)
				continue
			}

<<<<<<< HEAD
			result, err2 := stmt.Exec(nodeDetail.NodeUUID, nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads)
=======
			result, err2 := stmt.Exec(nodedetail.NodeUUID, nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads)
>>>>>>> f41ff24 (Refactoring packages structure)
			if err2 != nil {
				logger.Logger.Println(err2)
				_ = stmt.Close()
				continue
			}
			_ = stmt.Close()

<<<<<<< HEAD
			if config.Ipmi.Debug == "on" {
				result, err := result.LastInsertId()
				if err != nil {
					logger.Logger.Print("UpdateNodesDetail(): err=" + err.Error())
				} else {
					logger.Logger.Print("UpdateNodesDetail(): result=" + strconv.Itoa(int(result)))
				}
			}
=======
			logger.Logger.Println(result.LastInsertId())
>>>>>>> f41ff24 (Refactoring packages structure)
		} else {
			sql = "update node_detail set cpu_model = ?, cpu_processors = ?, cpu_threads = ? where node_uuid = ?"
			stmt, err := mysql.Db.Prepare(sql)
			if err != nil {
				logger.Logger.Println(err)
				continue
			}

<<<<<<< HEAD
			result, err2 := stmt.Exec(nodeDetail.CPUModel, nodeDetail.CPUProcessors, nodeDetail.CPUThreads, nodeDetail.NodeUUID)
=======
			result, err2 := stmt.Exec(nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads, nodedetail.NodeUUID)
>>>>>>> f41ff24 (Refactoring packages structure)
			if err2 != nil {
				logger.Logger.Println(err2)
				_ = stmt.Close()
				continue
			}
			_ = stmt.Close()

			if config.Ipmi.Debug == "on" {
<<<<<<< HEAD
				result, err := result.LastInsertId()
				if err != nil {
					logger.Logger.Print("UpdateNodesDetail(): err=" + err.Error())
				} else {
					logger.Logger.Print("UpdateNodesDetail(): result=" + strconv.Itoa(int(result)))
				}
			}
		}

		nodedetails = append(nodedetails, nodeDetail)
=======
				logger.Logger.Println(result.LastInsertId())
			}
		}

		nodedetails = append(nodedetails, nodedetail)
>>>>>>> f41ff24 (Refactoring packages structure)
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
