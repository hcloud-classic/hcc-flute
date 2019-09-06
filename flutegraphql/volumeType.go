package cellographql

import "github.com/graphql-go/graphql"

var volumeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Volume",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"server_uuid": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
