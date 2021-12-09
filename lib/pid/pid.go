package pid

import (
	"hcc/flute/lib/fileutil"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

var flutePIDFileLocation = "/var/run"
var flutePIDFile = "/var/run/flute.pid"

// IsFluteRunning : Check if flute is running
func IsFluteRunning() (running bool, pid int, err error) {
	if _, err := os.Stat(flutePIDFile); os.IsNotExist(err) {
		return false, 0, nil
	}

	pidStr, err := ioutil.ReadFile(flutePIDFile)
	if err != nil {
		return false, 0, err
	}

	flutePID, _ := strconv.Atoi(string(pidStr))

	proc, err := os.FindProcess(flutePID)
	if err != nil {
		return false, 0, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, flutePID, nil
	}

	return false, 0, nil
}

// WriteFlutePID : Write flute PID to "/var/run/flute.pid"
func WriteFlutePID() error {
	pid := os.Getpid()

	err := fileutil.CreateDirIfNotExist(flutePIDFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(flutePIDFile, strconv.Itoa(pid))
	if err != nil {
		return err
	}

	return nil
}

// DeleteFlutePID : Delete the flute PID file
func DeleteFlutePID() error {
	err := fileutil.DeleteFile(flutePIDFile)
	if err != nil {
		return err
	}

	return nil
}
