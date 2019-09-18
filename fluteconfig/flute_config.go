package fluteconfig

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/flute/flute.conf"

type fluteConfig struct {
	MysqlConfig *goconf.Section
	HttpConfig  *goconf.Section
}

var FluteConfig fluteConfig

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
##### CONFIG END #####
-----------------------------------*/