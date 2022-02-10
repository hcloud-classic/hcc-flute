package config

type rsakey struct {
	PrivateKeyFile string `goconf:"rsakey:private_key_file"` // PrivateKeyFile : RSA private key file for decrypt mysqld password
}

// Rsakey : rsakey config structure
var Rsakey rsakey
