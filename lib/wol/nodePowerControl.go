package wol

import (
	"github.com/linde12/gowol"
)

func OnNode(macAddr string) error {
	packet, err := gowol.NewMagicPacket(macAddr)
	if err != nil {
		return err
	}

	err = packet.SendPort("255.255.255.255", "7")
	err = packet.SendPort("255.255.255.255", "9")
	if err != nil {
		return err
	}

	return nil
}
