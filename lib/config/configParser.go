package config

import (
	"hcc/flute/lib/logger"
	"strings"

	"github.com/Terry-Mao/goconf"
)

var conf = goconf.New()
var config = fluteConfig{}
var err error

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		logger.Logger.Panicln("no mysql section")
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseHTTP() {
	config.HTTPConfig = conf.Get("http")
	if config.HTTPConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	HTTP = http{}
	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseIpmi() {
	config.IpmiConfig = conf.Get("ipmi")
	if config.IpmiConfig == nil {
		logger.Logger.Panicln("no ipmi section")
	}

	Ipmi = ipmi{}
	Ipmi.Debug, err = config.IpmiConfig.String("debug")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.BMCIPList, err = config.IpmiConfig.String("bmc_ip_list")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	Ipmi.BMCIPListArray = strings.Split(Ipmi.BMCIPList, ",")

	Ipmi.Username, err = config.IpmiConfig.String("username")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.Password, err = config.IpmiConfig.String("password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.RequestTimeoutMs, err = config.IpmiConfig.Int("request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.RequestRetry, err = config.IpmiConfig.Int("request_retry")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckAllIntervalMs, err = config.IpmiConfig.Int("check_all_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckStatusIntervalMs, err = config.IpmiConfig.Int("check_status_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckNodesDetailIntervalMs, err = config.IpmiConfig.Int("check_nodes_detail_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.BaseboardNICNoPXE, err = config.IpmiConfig.Int("baseboard_nic_no_pxe")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.BaseboardNICNoBMC, err = config.IpmiConfig.Int("baseboard_nic_no_bmc")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseWOL() {
	config.WOLConfig = conf.Get("wol")
	if config.WOLConfig == nil {
		logger.Logger.Panicln("no wol section")
	}

	WOL = wol{}
	WOL.BroadcastAddress, err = config.WOLConfig.String("broadcast_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

// Init : Parse config file and initialize config structure
func Init() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseMysql()
	parseHTTP()
	parseIpmi()
	parseWOL()
}
