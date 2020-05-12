package init

import (
	"errors"
	"hcc/flute/lib/logger"
)

func loggerInit() error {
	if !logger.Prepare() {
		return errors.New("error occurred while preparing logger")
	}

	return nil
}
