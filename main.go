package main

import (
	"fmt"
	"hcc/flute/driver/grpcsrv"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/syscheck"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		log.Fatalf("syscheck.CheckRoot(): %v", err.Error())
	}

	err = logger.Init()
	if err != nil {
		log.Fatalf("logger.Init(): %v", err.Error())
	}

	config.Init()

	err = mysql.Init()
	if err != nil {
		logger.Logger.Fatalf("mysql.Init(): %v", err.Error())
	}

	err = ipmi.BMCIPParser()
	if err != nil {
		logger.Logger.Fatalf("ipmi.BMCIPParser(): %v", err.Error())
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

	grpcsrv.Init()
}
