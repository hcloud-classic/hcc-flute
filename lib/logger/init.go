package logger

import "errors"

// Init : Prepare logger
func Init() error {
	if !prepare() {
		return errors.New("error occurred while preparing logger")
	}

	return nil
}

// End : Close logger
func End() {
	_ = FpLog.Close()
}
