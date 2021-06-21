package main

import (
	"fmt"
	"hcc/flute/action/grpc/client"
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
		errors.NewHccError(errors.FluteInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(logger.Logger)

	config.Init()

	err = mysql.Init()
	if err != nil {
		errors.NewHccError(errors.FluteInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
	}

	err = client.Init()
	if err != nil {
		errors.NewHccError(errors.FluteInternalInitFail, "client.Init(): "+err.Error()).Fatal()
	}

	err = ipmi.BMCIPParser()
	if err != nil {
		errors.NewHccError(errors.FluteInternalInitFail, "ipmi.BMCIPParser(): "+err.Error()).Fatal()
	}

	logger.Logger.Println("Starting ipmi.CheckNodeAll(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodeAllIntervalMs)) + "ms")
	ipmi.CheckNodeAll()
	logger.Logger.Println("Starting ipmi.CheckNodeStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodeStatusIntervalMs)) + "ms")
	ipmi.CheckNodeStatus()
	logger.Logger.Println("Starting ipmi.CheckServerStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckServerStatusIntervalMs)) + "ms")
	ipmi.CheckServerStatus()
	logger.Logger.Println("Starting ipmi.CheckNodeDetail(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodeDetailIntervalMs)) + "ms")
	ipmi.CheckNodeDetail()
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
