package main

import (
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
>>>>>>> 4afd3e80898e7f57c1dec709a37df8b08235a21b
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
	if !checkroot.CheckRoot() {
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

	logger.Logger.Println("Starting ipmi.CheckAll(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckAllIntervalMs)) + "ms")
	ipmi.CheckAll()
	logger.Logger.Println("Starting ipmi.CheckStatus(). Interval is " + strconv.Itoa(int(config.Ipmi.CheckStatusIntervalMs)) + "ms")
	ipmi.CheckStatus()

	http.Handle("/graphql", graphql.Handler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
>>>>>>> 50a1eafd315a248c5306efee36c4307de82b59cb
	}
}
