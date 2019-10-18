package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/flute/flute.conf"

type fluteConfig struct {
	MysqlConfig    *goconf.Section
	HTTPConfig     *goconf.Section
	RabbitMQConfig *goconf.Section
	IpmiConfig     *goconf.Section
}

/*-----------------------------------
         Config File Example

##### CONFIG START #####
[mysql]
id user
password pass
address 111.111.111.111
port 9999
database db_name

[http]
port 8888

[rabbitmq]
rabbitmq_id user
rabbitmq_password pass
rabbitmq_address 555.555.555.555
rabbitmq_port 15672

[ipmi]
debug off
bmc_ip_list 172.31.0.1,172.31.0.2,172.31.0.3,172.31.0.4
username user
password pass
request_timeout_ms 5000
request_retry 3
check_all_interval_ms 10000
check_status_interval_ms 5000
check_nodes_detail_interval_ms 15000
baseboard_nic_no_pxe 2
baseboard_nic_no_bmc 3
##### CONFIG END #####
-----------------------------------*/
