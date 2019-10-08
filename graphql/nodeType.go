package graphql

import "github.com/graphql-go/graphql"

var nodeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Node",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"bmc_mac_addr": &graphql.Field{
				Type: graphql.String,
			},
			"bmc_ip": &graphql.Field{
				Type: graphql.String,
			},
			"pxe_mac_addr": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"cpu_cores": &graphql.Field{
				Type: graphql.Int,
			},
			"memory": &graphql.Field{
				Type: graphql.Int,
			},
			"desc": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"active": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
