package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/flute/flute.conf"

type fluteConfig struct {
	RsakeyConfig  *goconf.Section
	MysqlConfig   *goconf.Section
	GrpcConfig    *goconf.Section
	IpmiConfig    *goconf.Section
	HornConfig    *goconf.Section
	ViolinConfig  *goconf.Section
	HarpConfig    *goconf.Section
	PiccoloConfig *goconf.Section
}

