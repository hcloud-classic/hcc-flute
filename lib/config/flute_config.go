package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/flute/flute.conf"

type fluteConfig struct {
	MysqlConfig *goconf.Section
	GrpcConfig  *goconf.Section
	IpmiConfig  *goconf.Section
}

