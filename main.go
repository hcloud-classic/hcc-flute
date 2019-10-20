package main

import (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"hcc/flute/action/graphql"
	fluteEnd "hcc/flute/end"
	fluteInit "hcc/flute/init"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
=======
	"hcloud-flute/checkroot"
	"hcloud-flute/config"
	"hcloud-flute/graphql"
	"hcloud-flute/logger"
	"hcloud-flute/mysql"
>>>>>>> 50a1eafd315a248c5306efee36c4307de82b59cb
=======
	"hcc/flute/checkroot"
	"hcc/flute/config"
	"hcc/flute/graphql"
	"hcc/flute/ipmi"
	"hcc/flute/logger"
	"hcc/flute/mysql"
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 4afd3e80898e7f57c1dec709a37df8b08235a21b
=======
	"hcc/flute/rabbitmq"
>>>>>>> 621c6e9 (ish + cgs work)
=======
	"hcc/flute/rabbitmq"
>>>>>>> 6df92f308119d4455f03636275411261b0c45f72
=======
	"hcc/flute/action/graphql"
	"hcc/flute/action/rabbitmq"
	"hcc/flute/lib/config"
	"hcc/flute/lib/ipmi"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/lib/syscheck"
>>>>>>> f41ff24 (Refactoring packages structure)
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
>>>>>>> 50a1eafd315a248c5306efee36c4307de82b59cb
	}
}
