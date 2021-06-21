package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/flute/logger"
)

// Parser : Parse config file
func Parser() {
	var conf = goconf.New()
	var config = fluteConfig{}
	var err error

	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

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

	config.HTTPConfig = conf.Get("http")
	if config.HTTPConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	HTTP = http{}
	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	config.IpmiConfig = conf.Get("ipmi")
	if config.IpmiConfig == nil {
		logger.Logger.Panicln("no ipmi section")
	}

	Ipmi = ipmi{}
	Ipmi.Username, err = config.IpmiConfig.String("username")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.Password, err = config.IpmiConfig.String("password")
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
}
