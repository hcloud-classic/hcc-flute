package config

type wol struct {
	BroadcastAddress string `goconf:"wol:broadcast_address"` // BroadcastAddress : Broadcast address for sending wol packet
}

// WOL : wol config structure
var WOL wol
