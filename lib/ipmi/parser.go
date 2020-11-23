package ipmi

import (
	"errors"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
<<<<<<< HEAD
	"hcc/flute/model"
=======
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
	"strings"
)

// BMCIPParser : Parse IP list of BMC and set active flags to database
func BMCIPParser() error {
	for _, ip := range config.Ipmi.BMCIPListArray {
		ipPart := strings.Split(ip, ".")
		if len(ipPart) != 4 {
			return errors.New("BMC IP list contains invalid IP address")
		}
	}

	sqlStr := "select bmc_ip from node"
	stmt, err := mysql.Db.Query(sqlStr)
	if err != nil {
		logger.Logger.Println(err)
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

<<<<<<< HEAD
	var nodes []model.Node
=======
	var nodes []types.Node
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
	var bmcIP string

	for stmt.Next() {
		err := stmt.Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return err
		}

<<<<<<< HEAD
		node := model.Node{BmcIP: bmcIP}
=======
		node := types.Node{BmcIP: bmcIP}
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
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
			logger.Logger.Println(err)
			return err
		}

		_, err2 := stmt.Exec(node.BmcIP)
		if err2 != nil {
			logger.Logger.Println(err2)
			return err2
		}

		_ = stmt.Close()
	}

	return nil
}

// BMCIPParserCheckActive : Check BMC IP list from config file and change active flag of given BMC IP from database
func BMCIPParserCheckActive(bmcIP string) error {
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
		logger.Logger.Println(err)
		return err
	}

	_, err2 := stmt.Exec(bmcIP)
	if err2 != nil {
		logger.Logger.Println(err2)
		return err2
	}

	_ = stmt.Close()

	return nil
}
