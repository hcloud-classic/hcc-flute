package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/flute/flute.conf"

type fluteConfig struct {
	MysqlConfig  *goconf.Section
	GrpcConfig   *goconf.Section
	IpmiConfig   *goconf.Section
	ViolinConfig *goconf.Section
}

/*-----------------------------------
         Config File Example
/*-----------------------------------
[mysql]
id root
password qwe1212!Q
address 192.168.110.240
port 3306
database flute

[grpc]
port 7000

[ipmi]
debug off
bmc_ip_list 172.31.0.10,172.31.0.1,172.31.0.3
username admin
password qwe1212!Q
request_timeout_ms 5000
request_retry 3
check_node_all_interval_ms 10000
check_node_status_interval_ms 1000
check_server_status_interval_ms 5000
check_node_detail_interval_ms 15000
baseboard_nic_num_pxe 2
baseboard_nic_num_bmc 3
-----------------------------------*/
