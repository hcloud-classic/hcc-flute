package main

import (
	"hcc/flute/action/graphql"
	"hcc/flute/lib/logger"
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

	config.Init()

	err = mysql.Init()
	if err != nil {
		panic(err)
	}
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
