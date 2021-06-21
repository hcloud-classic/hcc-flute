package iputil

import (
	"errors"
	"net"
)

func checkAClassPrivate(IP net.IP) bool {
	if IP[0] == 10 {
		return true
	}

	return false
}

func checkBClassPrivate(IP net.IP) bool {
	if IP[0] == 172 &&
		(IP[1] >= 16 && IP[1] <= 31) {
		return true
	}

	return false
}

func checkCClassPrivate(IP net.IP) bool {
	if IP[0] == 192 && IP[1] == 168 {
		return true
	}

	return false
}

// CheckPrivateSubnet : Check if given network address is private network address.
// Return error if given IP address is invalid or is not a network address.
// Return true if it is private address, return false otherwise.
func CheckPrivateSubnet(IP string, Netmask string) (bool, error) {
	netNetwork, err := CheckNetwork(IP, Netmask)
	if err != nil {
		return false, err
	}

	if netNetwork.IP.String() != IP {
		return false, errors.New("CheckPrivateSubnet(): invalid network address")
	}

	if checkAClassPrivate(netNetwork.IP) ||
		checkBClassPrivate(netNetwork.IP) ||
		checkCClassPrivate(netNetwork.IP) {
		return true, nil
	}

	return false, nil
}
