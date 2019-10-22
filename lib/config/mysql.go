package config

type mysql struct {
	ID       string `goconf:"mysql:id"`       // ID : MySQL server login id
	Password string `goconf:"mysql:password"` // Password : MySQL server login password
	Address  string `goconf:"mysql:address"`  // Address : MySQL server address
	Port     int64  `goconf:"mysql:port"`     // Port : MySQL server port number
	Database string `goconf:"mysql:database"` // Database : MySQL server database name of module
}

// Mysql : mysql config structure
var Mysql mysql
