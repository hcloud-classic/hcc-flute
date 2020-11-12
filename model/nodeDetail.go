package model

// NodeDetail - cgs
type NodeDetail struct {
	NodeUUID      string `json:"node_uuid"`
	CPUModel      string `json:"cpu_model"`
	CPUProcessors int    `json:"cpu_processors"`
	CPUThreads    int    `json:"cpu_threads"`
}
