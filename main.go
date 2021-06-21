package main

import (
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
	defer logger.FpLog.Close()

	config.ConfigParser()

	err := mysql.Prepare()
	if err != nil {
		return
	}
	defer mysql.Db.Close()

	http.Handle("/graphql", graphql.GraphqlHandler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.Http.Port)))
	err = http.ListenAndServe(":" + strconv.Itoa(int(config.Http.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
>>>>>>> 50a1eafd315a248c5306efee36c4307de82b59cb
	}
}
