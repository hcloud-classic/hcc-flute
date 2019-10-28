package wol

import (
	"github.com/linde12/gowol"
	"hcc/flute/lib/config"
)

func OnNode(macAddr string) error {
	packet, err := gowol.NewMagicPacket(macAddr)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		_ = packet.SendPort(config.WOL.BroadcastAddress, "7")
		_ = packet.SendPort(config.WOL.BroadcastAddress, "9")
	}

	return nil
}
