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

func delayMillisecond(n time.Duration) {
	time.Sleep(n * time.Millisecond)
}

func checkAllLock()  {
	checkAllLocked = true
}

func checkAllUnlock()  {
	checkAllLocked = false
}

func checkStatusLock()  {
	checkStatusLocked = true
}

func checkStatusUnlock()  {
	checkStatusLocked = false
}

func UpdateAllNodes()(interface{}, error) {
	var nodes []types.Node
	var uuid string
	var ipmiIP string

	sql := "select uuid, ipmi_ip from node"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &ipmiIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		serialNo, err := GetSerialNo(ipmiIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		mac, err := GetBMCNICMac(ipmiIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		powerState, err := GetPowerState(ipmiIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		processors, err := GetProcessors(ipmiIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}
		println(processors)

		cores, err := GetProcessorsCores(ipmiIP, serialNo, processors)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}
		println(cores)

		memory, err := GetTotalSystemMemory(ipmiIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		node := types.Node{
			UUID:    uuid,
			MacAddr: mac,
			IpmiIP:  ipmiIP,
			Status:  powerState,
			CPU:     cores,
			Memory:  memory,
		}

		sql := "update node set mac_addr = ?, status = ?, cpu = ?, memory = ? where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, nil
		}

		result, err2 := stmt.Exec(node.MacAddr, node.Status, node.CPU, node.Memory, node.UUID)
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

func UpdateStatusNodes()(interface{}, error) {
	var nodes []types.Node
	var uuid string
	var ipmiIP string

	sql := "select uuid, ipmi_ip from node"
	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, nil
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &ipmiIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		serialNo, err := GetSerialNo(ipmiIP)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		powerState, err := GetPowerState(ipmiIP, serialNo)
		if err != nil {
			logger.Logger.Println(err)
			return nil, nil
		}

		node := types.Node{
			UUID: uuid,
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

func queueCheckAll() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckAll(): Rerun CheckAll() after " + strconv.Itoa(int(config.Ipmi.CheckAllIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckAllIntervalMs))
		CheckAll();
	}()
}

func queueCheckStatus() {
	go func() {
		if config.Ipmi.Debug == "on" {
			logger.Logger.Println("queueCheckStatus(): Rerun CheckStatus() after " + strconv.Itoa(int(config.Ipmi.CheckStatusIntervalMs)) + "ms")
		}
		delayMillisecond(time.Duration(config.Ipmi.CheckStatusIntervalMs))
		CheckStatus();
	}()
}

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