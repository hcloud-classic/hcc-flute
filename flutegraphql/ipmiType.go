package flutegraphql

import "github.com/graphql-go/graphql"

var ipmiType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Ipmi",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"node_ip": &graphql.Field{
				Type: graphql.String,
			},
			"node_uuid": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
