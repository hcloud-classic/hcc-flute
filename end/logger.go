package end

import "hcc/flute/lib/logger"

func loggerEnd() {
	_ = logger.FpLog.Close()
}
