package main

import (
	"hcc/flute/action/graphql"
<<<<<<< HEAD
	fluteEnd "hcc/flute/end"
	fluteInit "hcc/flute/init"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
=======
	"hcc/flute/action/rabbitmq"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/syscheck"
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
	"net/http"
	"strconv"
)

<<<<<<< HEAD
func init() {
	err := fluteInit.MainInit()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		fluteEnd.MainEnd()
	}()

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
=======
func main() {
	if !syscheck.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := mysql.Prepare()
	if err != nil {
		return
	}
	defer func() {
		_ = mysql.Db.Close()
	}()

	err = ipmi.BMCIPParser()
	if err != nil {
		return
	}

	logger.Logger.Println("Starting ipmi.CheckAll(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckAllIntervalMs)) + "ms")
	ipmi.CheckAll()
	logger.Logger.Println("Starting ipmi.CheckStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckStatusIntervalMs)) + "ms")
	ipmi.CheckStatus()
	logger.Logger.Println("Starting ipmi.CheckNodesDetail(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckNodesDetailIntervalMs)) + "ms")
	ipmi.CheckNodesDetail()

	err = rabbitmq.PrepareChannel()
	if err != nil {
		logger.Logger.Panic(err)
	}
	defer func() {
		_ = rabbitmq.Channel.Close()
	}()
	defer func() {
		_ = rabbitmq.Connection.Close()
	}()

	err = rabbitmq.OnNode()
	if err != nil {
		logger.Logger.Panic(err)
	}

	http.Handle("/graphql", graphql.Handler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
	}
}
