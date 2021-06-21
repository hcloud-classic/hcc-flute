package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/flute/flute.conf"

type fluteConfig struct {
	MysqlConfig *goconf.Section
	HttpConfig  *goconf.Section
}
