package config

type ipmi struct {
	Debug                           string   `goconf:"ipmi:debug"`       // Debug : Enable debug logs for IPMI
	BMCIPList                       string   `goconf:"ipmi:bmc_ip_list"` // BMCIPList : List of BMC IPs
	BMCIPListArray                  []string // BMCIPList : Array list of BMC IPs
	PasswordEncryptSecretKey        string   `goconf:"ipmi:password_encrypt_secret_key"`          // PasswordEncryptSecretKey : Secret key for encrypt IPMI user's password
	RequestTimeoutMs                int64    `goconf:"ipmi:request_timeout_ms"`                   // RequestTimeoutMs : Timeout for IPMI request
	RequestRetry                    int64    `goconf:"ipmi:request_retry"`                        // RequestRetry : Retry count for IPMI request
	CheckNodeAllIntervalMs          int64    `goconf:"ipmi:check_node_all_interval_ms"`           // CheckNodeAllIntervalMs : IPMI check interval for node's all infos (ms)
	CheckNodeStatusIntervalMs       int64    `goconf:"ipmi:check_node_status_interval_ms"`        // CheckNodeStatusIntervalMs : IPMI check interval for node's status (ms)
	CheckServerStatusIntervalMs     int64    `goconf:"ipmi:check_server_status_interval_ms"`      // CheckServerStatusIntervalMs : IPMI check interval for server's status (ms)
	UpdateNodeDetailRetryIntervalMs int64    `goconf:"ipmi:update_node_detail_retry_interval_ms"` // UpdateNodeDetailRetryIntervalMs : Node update retry interval by IPMI (ms)
	UpdateNodeUptimeIntervalMs      int64    `goconf:"ipmi:update_node_uptime_interval_ms"`       // UpdateNodeUptimeIntervalMs : Node uptime update interval (ms)
	BaseboardNICNumPXE              int64    `goconf:"ipmi:baseboard_nic_num_pxe"`                // BaseboardNICNoPXE : Baseboard NIC num used for PXE boot
	BaseboardNICNumBMC              int64    `goconf:"ipmi:baseboard_nic_num_bmc"`                // BaseboardNICNoIPMI : Baseboard NIC num used for control IPMI
}

// Ipmi : ipmi config structure
var Ipmi ipmi
