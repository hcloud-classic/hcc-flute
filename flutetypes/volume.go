package cellotypes

// Volume : Struct of volume
type Volume struct {
	UUID       string `json:"uuid"`
	Size       int    `json:"size"`
	Type       string `json:"type"`
	ServerUUID string `json:"server_uuid"`
}

// Volumes : Array struct of volumes
type Volumes struct {
	Volumes []Volume `json:"server"`
}
