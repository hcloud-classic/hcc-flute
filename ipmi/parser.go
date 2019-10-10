package ipmi

import (
	"hcc/flute/config"
	"hcc/flute/logger"
	"hcc/flute/mysql"
	"hcc/flute/types"
	"strings"
)

// BMCIPParser : Parse IP list of BMC and set active flags to database
func BMCIPParser() {
	for _, ip := range config.Ipmi.BMCIPListArray {
		ipPart := strings.Split(ip, ".")
		if len(ipPart) != 4 {
			logger.Logger.Panic("BMC IP list contains invalid IP address")
		}
	}

	sqlStr := "select bmc_ip from node"
	stmt, err := mysql.Db.Query(sqlStr)
	if err != nil {
		logger.Logger.Panic(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var nodes []types.Node
	var bmcIP string

	for stmt.Next() {
		err := stmt.Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
		}

		node := types.Node{BmcIP: bmcIP}
		nodes = append(nodes, node)
	}

	for _, node := range nodes {
		var ipMatched = false
		for _, ip := range config.Ipmi.BMCIPListArray {
			if node.BmcIP == ip {
				ipMatched = true
				break
			}
		}

		var sqlStr string
		if ipMatched {
			sqlStr = "update node set active = 1 where bmc_ip = ?"
		} else {
			sqlStr = "update node set active = 0 where bmc_ip = ?"
		}

		stmt, err := mysql.Db.Prepare(sqlStr)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		_, err2 := stmt.Exec(node.BmcIP)
		if err2 != nil {
			logger.Logger.Panic(err2)
		}

		_ = stmt.Close()
	}
}

func BMCIPParserCheckActive(bmcIP string) {
	var ipMatched = false
	for _, ip := range config.Ipmi.BMCIPListArray {
		if bmcIP == ip {
			ipMatched = true
			break
		}
	}

	var sqlStr string
	if ipMatched {
		sqlStr = "update node set active = 1 where bmc_ip = ?"
	} else {
		sqlStr = "update node set active = 0 where bmc_ip = ?"
	}

	stmt, err := mysql.Db.Prepare(sqlStr)
	if err != nil {
		logger.Logger.Panic(err.Error())
	}

	_, err2 := stmt.Exec(bmcIP)
	if err2 != nil {
		logger.Logger.Panic(err2)
	}

	_ = stmt.Close()
}
