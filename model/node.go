package model

import "time"

// Node : Struct of node
type Node struct {
	UUID        string    `json:"uuid"`
	BmcMacAddr  string    `json:"bmc_mac_addr"`
	BmcIP       string    `json:"bmc_ip"`
	PXEMacAddr  string    `json:"pxe_mac_addr"`
	Status      string    `json:"status"`
	CPUCores    int       `json:"cpu_cores"`
	Memory      int       `json:"memory"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Active      int       `json:"active"`
	ForceOff    bool      `json:"force_off"`
}

// Nodes : Array struct of nodes
type Nodes struct {
	Nodes []Node `json:"node"`
}

// NodeNum : Struct of number of nodes
type NodeNum struct {
	Number int `json:"number"`
}
