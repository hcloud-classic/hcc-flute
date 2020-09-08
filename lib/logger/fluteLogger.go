package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// LogName : Log folder name. Also used as log prefix.
var LogName = "flute"

// Logger : Pointer of logger
var Logger *log.Logger

// FpLog : File pointer of logger
var FpLog *os.File

// CreateDirIfNotExist : Make directory if not exist
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// Init : Initialize logger
func Init() error {
	// Create directory if not exist
	if _, err := os.Stat("/var/log/" + LogName); os.IsNotExist(err) {
		err = CreateDirIfNotExist("/var/log/" + LogName)
		if err != nil {
			return err
		}
	}

	now := time.Now()

	year := fmt.Sprintf("%d", now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())

	date := year + month + day

	FpLog, err := os.OpenFile("/var/log/"+LogName+"/"+
		LogName+"_"+date+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Logger = log.New(io.MultiWriter(FpLog, os.Stdout), LogName+"_logger: ", log.Ldate|log.Ltime)

	return nil
}

// End : Close logger
func End() {
	_ = FpLog.Close()
}

