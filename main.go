package main

import (
	"fmt"
	"hcc/flute/action/grpc/client"
	"hcc/flute/action/grpc/server"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/pid"
	"innogrid.com/hcloud-classic/hcc_errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	fluteRunning, flutePID, err := pid.IsFluteRunning()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if fluteRunning {
		fmt.Println("flute is already running. (PID: " + strconv.Itoa(flutePID) + ")")
		os.Exit(1)
	}
	err = pid.WriteFlutePID()
	if err != nil {
		_ = pid.DeleteFlutePID()
		fmt.Println(err)
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteFlutePID()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	config.Init()

	err = mysql.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteFlutePID()
	}

	err = client.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "client.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteFlutePID()
	}

	err = ipmi.BMCIPParser()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.FluteInternalInitFail, "ipmi.BMCIPParser(): "+err.Error()).Fatal()
		_ = pid.DeleteFlutePID()
	}

	logger.Logger.Println("Starting ipmi.CheckNodeAll(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodeAllIntervalMs)) + "ms")
	ipmi.CheckNodeAll()
	logger.Logger.Println("Starting ipmi.CheckNodeStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodeStatusIntervalMs)) + "ms")
	ipmi.CheckNodeStatus()
	logger.Logger.Println("Starting ipmi.CheckServerStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckServerStatusIntervalMs)) + "ms")
	ipmi.CheckServerStatus()
	logger.Logger.Println("Starting ipmi.UpdateNodeUptime(). Interval is " + strconv.Itoa(int(config.Ipmi.UpdateNodeUptimeIntervalMs)) + "ms")
	ipmi.UpdateNodeUptime(time.Now())
}

func end() {
	client.End()
	mysql.End()
	logger.End()
	_ = pid.DeleteFlutePID()
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
