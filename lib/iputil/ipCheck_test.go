package iputil

import (
	"net"
	"testing"
)

func Test_CheckIP(t *testing.T) {
	netIP := CheckValidIP("192.168.100.0")
	if netIP == nil {
		t.Fatal("wrong network IP")
	}

	mask, err := CheckNetmask("255.255.255.0")
	if err != nil {
		t.Fatal(err)
	}

	ipNet := net.IPNet{
		IP:   netIP,
		Mask: mask,
	}

	err = CheckIPisInSubnet(ipNet, "192.168.100.1")
	if err != nil {
		t.Fatal(err)
	}
}
