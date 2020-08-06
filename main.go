package main

import (
	"hcc/flute/action/graphql"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/syscheck"
	"net/http"
	"strconv"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		panic(err)
	}

	config.Parser()

	err = mysql.Init()
	if err != nil {
		panic(err)
	}

	err = ipmi.BMCIPParser()
	if err != nil {
		panic(err)
	}

	logger.Logger.Println("Starting ipmi.CheckAll(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckAllIntervalMs)) + "ms")
	ipmi.CheckAll()
	logger.Logger.Println("Starting ipmi.CheckStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckStatusIntervalMs)) + "ms")
	ipmi.CheckStatus()
	logger.Logger.Println("Starting ipmi.CheckNodesDetail(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodesDetailIntervalMs)) + "ms")
	ipmi.CheckNodesDetail()
}

func main() {
	defer func() {
		mysql.End()
		logger.End()
	}()

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
	}
}
