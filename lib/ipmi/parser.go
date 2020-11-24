package ipmi

import (
	pb "hcc/flute/action/grpc/pb/rpcflute"
	"hcc/flute/lib/config"
	"hcc/flute/lib/iputil"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"net"
)

// BMCIPParser : Parse IP list of BMC and set active flags to database
func BMCIPParser() error {
	logger.Logger.Println("Parsing 'bmc_ip_list' from 'flute.conf'...")

	for _, cidr := range config.Ipmi.BMCIPListArray {
		err := iputil.CheckCIDRStr(cidr)
		if err != nil {
			return err
		}

		_, _, err = net.ParseCIDR(cidr)
		if err != nil {
			return err
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

	var nodes []pb.Node
	var bmcIP string

	for stmt.Next() {
		err := stmt.Scan(&bmcIP)
		if err != nil {
			logger.Logger.Println(err)
			return err
		}

		nodes = append(nodes, pb.Node{BmcIP: bmcIP})
	}

	for i := range nodes {
		var ipMatched = false
		for _, ip := range config.Ipmi.BMCIPListArray {
			if nodes[i].BmcIP == ip {
				ipMatched = true
				break
			}
		}

		var sqlStr string
		if ipMatched {
			sqlStr = "update node set available = 1 where bmc_ip = ?"
		} else {
			sqlStr = "update node set available = 0 where bmc_ip = ?"
		}

		stmt, err := mysql.Db.Prepare(sqlStr)
		if err != nil {
			logger.Logger.Println(err)
			return err
		}

		_, err2 := stmt.Exec(nodes[i].BmcIP)
		if err2 != nil {
			logger.Logger.Println(err2)
			return err2
		}

		_ = stmt.Close()
	}

	return nil
}
