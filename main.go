package main

import (
	"fmt"
	"hcc/flute/action/grpc/server"
	"hcc/flute/lib/config"
	"hcc/flute/lib/errors"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	err := logger.Init()
	if err != nil {
		errors.SetErrLogger(logger.Logger)
		errors.NewHccError(errors.HarpInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(logger.Logger)

	config.Init()

	err = mysql.Init()
	if err != nil {
		errors.NewHccError(errors.HarpInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
	}

	err = ipmi.BMCIPParser()
	if err != nil {
		errors.NewHccError(errors.HarpInternalInitFail, "ipmi.BMCIPParser(): "+err.Error()).Fatal()
	}

	logger.Logger.Println("Starting ipmi.CheckAll(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckAllIntervalMs)) + "ms")
	ipmi.CheckAll()
	logger.Logger.Println("Starting ipmi.CheckStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckStatusIntervalMs)) + "ms")
	ipmi.CheckStatus()
	logger.Logger.Println("Starting ipmi.CheckNodesDetail(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodesDetailIntervalMs)) + "ms")
	ipmi.CheckNodesDetail()
}

func end() {
	mysql.End()
	logger.End()
}

func main() {
	// Catch the exit signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		end()
		fmt.Println("Exiting flute module...")
		os.Exit(0)
	}()

	server.Init()
}
