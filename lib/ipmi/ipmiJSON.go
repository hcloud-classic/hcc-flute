package ipmi

type status struct {
	State        string `json:"state"`
	Health       string `json:"health"`
	HealthRollup string `json:"health_rollup"`
}

type memory struct {
	ID            string `json:"id"`
	CapacityMB    int    `json:"capacity_mb"`
	Manufacture   string `json:"manufacture"`
	SerialNumber  string `json:"serial_number"`
	PartNumber    string `json:"part_number"`
	DeviceLocator string `json:"device_locator"`
	SpeedMhz      int    `json:"speed_mhz"`
	Status        status `json:"status"`
}

type cpu struct {
	ID          string `json:"id"`
	Socket      string `json:"socket"`
	Manufacture string `json:"manufacture"`
	Model       string `json:"model"`
	MaxSpeedMHz int    `json:"max_speed_mhz"`
	Cores       int    `json:"cores"`
	Threads     int    `json:"threads"`
	Status      status `json:"status"`
}

type nodeDetailData struct {
	Memories []memory `json:"memories"`
	CPUs     []cpu    `json:"cpus"`
}
