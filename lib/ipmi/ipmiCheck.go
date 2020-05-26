package ipmi

import (
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
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
	var nodes []model.Node
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

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		uuid, err := GetUUID(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		bmcMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoBMC), true)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		pxeMAC, err := GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoPXE), false)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		processors, err := GetProcessors(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		cpuCores, err := GetProcessorsCores(bmcIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		memory, err := GetTotalSystemMemory(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		node := model.Node{
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
			logger.Logger.Println(err)
			continue
		}

		result, err2 := stmt.Exec(node.BmcMacAddr, node.PXEMacAddr, node.CPUCores, node.Memory, node.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			_ = stmt.Close()
			continue
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
	var nodes []model.Node
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
			continue
		}

		serialNo, err := GetSerialNo(bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		powerState, err := GetPowerState(bmcIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			continue
		}

		node := model.Node{
			UUID:   uuid,
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
			logger.Logger.Println(result.LastInsertId())
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
