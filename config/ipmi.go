package config

type ipmi struct {
	Debug                 string `goconf:"ipmi:debug"`                    // Debug : Enable debug logs for IPMI
	Username              string `goconf:"ipmi:username"`                 // Username : IPMI http basicauth username
	Password              string `goconf:"ipmi:password"`                 // Password : IPMI http basicauth password
	CheckAllIntervalMs    int64  `goconf:"ipmi:check_all_interval_ms"`    // CheckAllIntervalMs : IMPI check interval for all infos (ms)
	CheckStatusIntervalMs int64  `goconf:"ipmi:check_status_interval_ms"` // CheckStatusIntervalMs : IMPI check interval for status (ms)
}

// Ipmi : ipmi config structure
var Ipmi ipmi
