package config

type ipmi struct {
	Debug                       string   `goconf:"ipmi:debug"`       // Debug : Enable debug logs for IPMI
	BMCIPList                   string   `goconf:"ipmi:bmc_ip_list"` // BMCIPList : List of BMC IPs
	BMCIPListArray              []string // BMCIPList : Array list of BMC IPs
	Username                    string   `goconf:"ipmi:username"`                        // Username : IPMI http basicauth username
	Password                    string   `goconf:"ipmi:password"`                        // Password : IPMI http basicauth password
	RequestTimeoutMs            int64    `goconf:"ipmi:request_timeout_ms"`              // RequestTimeoutMs : Timeout for IPMI request
	RequestRetry                int64    `goconf:"ipmi:request_retry"`                   // RequestRetry : Retry count for IPMI request
	CheckNodeAllIntervalMs      int64    `goconf:"ipmi:check_node_all_interval_ms"`      // CheckNodeAllIntervalMs : IPMI check interval for node's all infos (ms)
	CheckNodeStatusIntervalMs   int64    `goconf:"ipmi:check_node_status_interval_ms"`   // CheckNodeStatusIntervalMs : IPMI check interval for node's status (ms)
	CheckServerStatusIntervalMs int64    `goconf:"ipmi:check_server_status_interval_ms"` // CheckServerStatusIntervalMs : IPMI check interval for server's status (ms)
	CheckNodeDetailIntervalMs   int64    `goconf:"ipmi:check_node_detail_interval_ms"`   // CheckNodeDetailIntervalMs : IPMI check interval for node's detail (ms)
	BaseboardNICNumPXE          int64    `goconf:"ipmi:baseboard_nic_num_pxe"`           // BaseboardNICNoPXE : Baseboard NIC num used for PXE boot
	BaseboardNICNumBMC          int64    `goconf:"ipmi:baseboard_nic_num_bmc"`           // BaseboardNICNoIPMI : Baseboard NIC num used for control IPMI
}

// Ipmi : ipmi config structure
var Ipmi ipmi
