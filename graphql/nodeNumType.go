package graphql

import "github.com/graphql-go/graphql"

var nodeNum = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "NodeNum",
		Fields: graphql.Fields{
			"number": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
