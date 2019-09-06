package cellocheckroot

import (
	"fmt"
	"os"
)

// CheckRoot : Check root permission (Check if uid is 0)
func CheckRoot() bool {
	if os.Geteuid() != 0 {
		fmt.Println("Please run as root!")
		return false
	}

	return true
}
