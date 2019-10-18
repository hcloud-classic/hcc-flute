package types

import "time"

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
}

type Nodes struct {
	Nodes []Node `json:"node"`
}

type NodeNum struct {
	Number int `json:"number"`
}
