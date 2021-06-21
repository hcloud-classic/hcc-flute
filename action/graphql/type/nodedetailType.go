package graphqlType

import "github.com/graphql-go/graphql"

// NodeDetailType : Graphql object type of NodeDetail
var NodeDetailType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "NodeDetail",
		Fields: graphql.Fields{
			"node_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"cpu_model": &graphql.Field{
				Type: graphql.String,
			},
			"cpu_processors": &graphql.Field{
				Type: graphql.Int,
			},
			"cpu_threads": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
