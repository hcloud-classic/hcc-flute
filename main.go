package main

import (
	"GraphQL_Flute/flutecheckroot"
	"GraphQL_Flute/fluteconfig"
	"GraphQL_Flute/flutegraphql"
	"GraphQL_Flute/flutelogger"
	"GraphQL_Flute/flutemysql"
	"net/http"
	"strconv"
)

func main() {
	if !flutecheckroot.CheckRoot() {
		return
	}

	if !flutelogger.Prepare() {
		return
	}
	defer flutelogger.FpLog.Close()

	fluteconfig.ConfigParser()

	err := flutemysql.Prepare()
	if err != nil {
		return
	}
	defer flutemysql.Db.Close()

	http.Handle("/graphql", flutegraphql.GraphqlHandler)

	flutelogger.Logger.Println("Server is running on port " + strconv.Itoa(int(fluteconfig.Http.Port)))
	err = http.ListenAndServe(":" + strconv.Itoa(int(fluteconfig.Http.Port)), nil)
	if err != nil {
		flutelogger.Logger.Println("Failed to prepare http server!")
	}
}
