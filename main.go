package main

import (
	"hcc/flute/action/graphql"
	fluteEnd "hcc/flute/end"
	fluteInit "hcc/flute/init"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"net/http"
	"strconv"
)

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
	}
}
