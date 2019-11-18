package syscheck

import (
	"errors"
	"os"
)

// CheckRoot : Check root permission (Check if uid is 0)
func CheckRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("please run as root")
	}

	return nil
}
