package config

import (
	"github.com/Terry-Mao/goconf"
	"hcloud-flute/logger"
)

var Conf = goconf.New()

func Parser() {
	var config = fluteConfig{}
	var err error

	if err = Conf.Parse(configLocation)
	err != nil {
		logger.Logger.Panicln(err)
	}

	config.MysqlConfig = Conf.Get("mysql")
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

	config.HttpConfig = Conf.Get("http")
	if config.HttpConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	Http = http{}
	Http.Port, err = config.HttpConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}