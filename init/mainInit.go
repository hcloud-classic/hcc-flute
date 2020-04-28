package init

import (
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
)

// MainInit : Main initialization function
func MainInit() error {
	err := syscheckInit()
	if err != nil {
		return err
	}

	err = loggerInit()
	if err != nil {
		return err
	}

	config.Parser()

	err = mysqlInit()
	if err != nil {
		return err
	}

	// TODO : Temporary not using IPMI
	err = ipmi.BMCIPParser()
	if err != nil {
		return err
	}

	return nil
}
