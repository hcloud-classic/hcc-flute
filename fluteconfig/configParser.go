package fluteconfig

import (
	"GraphQL_Flute/flutelogger"
	"github.com/Terry-Mao/goconf"
)

var Conf = goconf.New()

func ConfigParser() {
	FluteConfig = fluteConfig{}
	var err error

	if err = Conf.Parse(configLocation)
	err != nil {
		flutelogger.Logger.Panicln(err)
	}

	FluteConfig.MysqlConfig = Conf.Get("mysql")
	if FluteConfig.MysqlConfig == nil {
		flutelogger.Logger.Panicln("no mysql section")
	}

	Mysql = mysql{}
	Mysql.ID, err = FluteConfig.MysqlConfig.String("id")
	if err != nil {
		flutelogger.Logger.Panicln(err)
	}

	Mysql.Password, err = FluteConfig.MysqlConfig.String("password")
	if err != nil {
		flutelogger.Logger.Panicln(err)
	}

	Mysql.Address, err = FluteConfig.MysqlConfig.String("address")
	if err != nil {
		flutelogger.Logger.Panicln(err)
	}

	Mysql.Port, err = FluteConfig.MysqlConfig.Int("port")
	if err != nil {
		flutelogger.Logger.Panicln(err)
	}

	Mysql.Database, err = FluteConfig.MysqlConfig.String("database")
	if err != nil {
		flutelogger.Logger.Panicln(err)
	}

	FluteConfig.HttpConfig = Conf.Get("http")
	if FluteConfig.HttpConfig == nil {
		flutelogger.Logger.Panicln("no http section")
	}

	Http = http{}
	Http.Port, err = FluteConfig.HttpConfig.Int("port")
	if err != nil {
		flutelogger.Logger.Panicln(err)
	}
}