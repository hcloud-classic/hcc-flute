package types

import "time"

// Node : Struct of node
type Node struct {
	UUID      string `json:"uuid"`
	MacAddr   string `json:"mac_addr"`
	IpmiIP    string `json:"ipmi_ip"`
	Status    string `json:"status"`
	Cpu       int `json:"cpu"`
	Memory    int `json:"memory"`
	Detail    string `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

// Nodes : Array struct of nodes
type Nodes struct {
	Nodes []Node `json:"node"`
}
