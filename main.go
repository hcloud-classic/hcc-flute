package main

import (
	"hcloud-flute/checkroot"
	"hcloud-flute/config"
	"hcloud-flute/graphql"
	"hcloud-flute/logger"
	"hcloud-flute/mysql"
	"net/http"
	"strconv"
)

func main() {
	if !checkroot.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer logger.FpLog.Close()

	config.Parser()

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
	}
}
