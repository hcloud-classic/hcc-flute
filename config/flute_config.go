package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/flute/flute.conf"

type fluteConfig struct {
	MysqlConfig *goconf.Section
	HTTPConfig  *goconf.Section
	IpmiConfig  *goconf.Section
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

[ipmi]
debug off
username user
password pass
check_all_interval_ms 10000
check_status_interval_ms 1000
##### CONFIG END #####
-----------------------------------*/
