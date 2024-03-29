package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/flute/lib/logger"
	"strings"
)

var conf = goconf.New()
var config = fluteConfig{}
var err error

func parseRsakey() {
	config.RsakeyConfig = conf.Get("rsakey")
	if config.RsakeyConfig == nil {
		logger.Logger.Panicln("no rsakey section")
	}

	Rsakey.PrivateKeyFile, err = config.RsakeyConfig.String("private_key_file")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

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

	Ipmi.PasswordEncryptSecretKey, err = config.IpmiConfig.String("password_encrypt_secret_key")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	if Ipmi.PasswordEncryptSecretKey == "" {
		logger.Logger.Panicln("password_encrypt_secret_key should not be empty value")
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

	Ipmi.UpdateNodeDetailRetryIntervalMs, err = config.IpmiConfig.Int("update_node_detail_retry_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.UpdateNodeUptimeIntervalMs, err = config.IpmiConfig.Int("update_node_uptime_interval_ms")
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

	Ipmi.CheckNodeOffConfirmIntervalMs, err = config.IpmiConfig.Int("check_node_off_confirm_interval_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.CheckNodeOffConfirmRetryCounts, err = config.IpmiConfig.Int("check_node_off_confirm_retry_counts")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.ServerStatusCheckPowerOnTimeOutSec, err = config.IpmiConfig.Int("server_status_check_power_on_timeout_sec")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.ServerStatusCheckBootingTimeoutSec, err = config.IpmiConfig.Int("server_status_check_booting_timeout_sec")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.ServerStatusCheckSSHPort, err = config.IpmiConfig.Int("server_status_check_ssh_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Ipmi.ServerStatusCheckVNCPort, err = config.IpmiConfig.Int("server_status_check_vnc_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseHorn() {
	config.HornConfig = conf.Get("horn")
	if config.HornConfig == nil {
		logger.Logger.Panicln("no horn section")
	}

	Horn = horn{}
	Horn.ServerAddress, err = config.HornConfig.String("horn_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.ServerPort, err = config.HornConfig.Int("horn_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.ConnectionTimeOutMs, err = config.HornConfig.Int("horn_connection_timeout_ms")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.ConnectionRetryCount, err = config.HornConfig.Int("horn_connection_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Horn.RequestTimeoutMs, err = config.HornConfig.Int("horn_request_timeout_ms")
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

func parseHarp() {
	config.HarpConfig = conf.Get("harp")
	if config.HarpConfig == nil {
		logger.Logger.Panicln("no harp section")
	}

	Harp = harp{}
	Harp.ServerAddress, err = config.HarpConfig.String("harp_server_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Harp.ServerPort, err = config.HarpConfig.Int("harp_server_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Harp.RequestTimeoutMs, err = config.HarpConfig.Int("harp_request_timeout_ms")
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

	parseRsakey()
	parseMysql()
	parseGrpc()
	parseIpmi()
	parseHorn()
	parseViolin()
	parseHarp()
	parsePiccolo()
}
