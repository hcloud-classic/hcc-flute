package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/flute/lib/logger"
	"strings"
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

	Mysql.ConnectionRetryCount, err = config.MysqlConfig.Int("connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.ConnectionRetryIntervalMs, err = config.MysqlConfig.Int("connection_retry_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		logger.Logger.Panicln("no grpc section")
	}

	Grpc.Port, err = config.GrpcConfig.Int("port")
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

	Ipmi.CheckNodeAllIntervalMs, err = config.IpmiConfig.Int("check_node_all_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckNodeStatusIntervalMs, err = config.IpmiConfig.Int("check_node_status_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckServerStatusIntervalMs, err = config.IpmiConfig.Int("check_server_status_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckNodeDetailIntervalMs, err = config.IpmiConfig.Int("check_node_detail_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.BaseboardNICNumPXE, err = config.IpmiConfig.Int("baseboard_nic_num_pxe")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.BaseboardNICNumBMC, err = config.IpmiConfig.Int("baseboard_nic_num_bmc")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseViolin() {
	config.ViolinConfig = conf.Get("violin")
	if config.ViolinConfig == nil {
		logger.Logger.Panicln("no violin section")
	}

	Violin = violin{}
	Violin.ServerAddress, err = config.ViolinConfig.String("violin_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.ServerPort, err = config.ViolinConfig.Int("violin_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.ConnectionTimeOutMs, err = config.ViolinConfig.Int("violin_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.ConnectionRetryCount, err = config.ViolinConfig.Int("violin_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Violin.RequestTimeoutMs, err = config.ViolinConfig.Int("violin_request_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parsePiccolo() {
	config.PiccoloConfig = conf.Get("piccolo")
	if config.PiccoloConfig == nil {
		logger.Logger.Panicln("no piccolo section")
	}

	Piccolo = piccolo{}
	Piccolo.ServerAddress, err = config.PiccoloConfig.String("piccolo_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.ServerPort, err = config.PiccoloConfig.Int("piccolo_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.ConnectionTimeOutMs, err = config.PiccoloConfig.Int("piccolo_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.ConnectionRetryCount, err = config.PiccoloConfig.Int("piccolo_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Piccolo.RequestTimeoutMs, err = config.PiccoloConfig.Int("piccolo_request_timeout_ms")
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
	parseGrpc()
	parseIpmi()
	parseViolin()
	parsePiccolo()
}
