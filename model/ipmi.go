package model

// Ipmi : Struct of ipmi
type Ipmi struct {
	UUID string `json:"uuid"`
}

// Ipmis : Array struct of ipmis
type Ipmis struct {
	Ipmis []Ipmi `json:"ipmi"`
}
