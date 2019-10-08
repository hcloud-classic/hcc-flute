package types

// NodeDetail : Struct of node_detail
type NodeDetail struct {
	NodeUUID      string `json:"node_uuid"`
	CPUModel      string `json:"cpu_model"`
	CPUProcessors int    `json:"cpu_processors"`
	CPUThreads    int    `json:"cpu_threads"`
}

// NodeDetails : Array struct of node_details
type NodeDetails struct {
	Nodes []Node `json:"node_detail"`
}
