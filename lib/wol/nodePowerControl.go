package wol

import (
	"github.com/linde12/gowol"
)

func OnNode(macAddr string) error {
	packet, err := gowol.NewMagicPacket(macAddr)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		_ = packet.SendPort("255.255.255.255", "7")
		_ = packet.SendPort("255.255.255.255", "9")
	}

	return nil
}
