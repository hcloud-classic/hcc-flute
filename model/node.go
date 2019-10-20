package model

import "time"

<<<<<<< HEAD
// Node : Contain infos of a node
type Node struct {
	UUID        string    `json:"uuid"`
	ServerUUID  string    `json:"server_uuid"`
=======
// Node : Struct of node
type Node struct {
	UUID        string    `json:"uuid"`
>>>>>>> f41ff24 (Refactoring packages structure)
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

<<<<<<< HEAD
// Nodes : Contain a node list
=======
// Nodes : Array struct of nodes
>>>>>>> f41ff24 (Refactoring packages structure)
type Nodes struct {
	Nodes []Node `json:"node"`
}

<<<<<<< HEAD
// NodeNum : Contain the number of nodes
=======
// NodeNum : Struct of number of nodes
>>>>>>> f41ff24 (Refactoring packages structure)
type NodeNum struct {
	Number int `json:"number"`
}
