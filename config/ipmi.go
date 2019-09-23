package config

type ipmi struct {
	Username string `goconf:"ipmi:username"` // Username : IPMI http basicauth username
	Password string `goconf:"ipmi:password"` // Password : IPMI http basicauth password
}

// Ipmi : ipmi config structure
var Ipmi ipmi
