package types

import "time"

// Node : Struct of node
type Node struct {
	UUID       string    `json:"uuid"`
	BmcMacAddr string    `json:"bmc_mac_addr"`
	BmcIP      string    `json:"bmc_ip"`
	PXEMacAddr string    `json:"pxe_mac_addr"`
	Status     string    `json:"status"`
	CPUCores   int       `json:"cpu_cores"`
	Memory     int       `json:"memory"`
	Desc       string    `json:"desc"`
	CreatedAt  time.Time `json:"created_at"`
	Active     int       `json:"active"`
}

// Nodes : Array struct of nodes
type Nodes struct {
	Nodes []Node `json:"node"`
}
