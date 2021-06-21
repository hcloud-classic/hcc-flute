package graphql

import (
	"github.com/graphql-go/graphql"
	graphqlType "hcc/flute/action/graphql/type"
	"hcc/flute/dao"
	"hcc/flute/lib/logger"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// node DB
			"node": &graphql.Field{
				Type:        graphqlType.NodeType,
				Description: "Get a node by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node")
					return dao.ReadNode(params.Args)
				},
			},
			"list_node": &graphql.Field{
				Type:        graphql.NewList(graphqlType.NodeType),
				Description: "Get node list",
				Args: graphql.FieldConfigArgument{
					"server_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"bmc_mac_addr": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"bmc_ip": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"pxe_mac_addr": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"cpu_cores": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"memory": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"active": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_node")
					return dao.ReadNodeList(params.Args)
				},
			},
			"all_node": &graphql.Field{
				Type:        graphql.NewList(graphqlType.NodeType),
				Description: "Get all node list",
				Args: graphql.FieldConfigArgument{
					"active": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_node")
					return dao.ReadNodeAll(params.Args)
				},
			},
			"num_node": &graphql.Field{
				Type:        graphqlType.NodeNumType,
				Description: "Get the number of node",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: num_node")
					return dao.ReadNodeNum(params.Args)
				},
			},
			// detail_node DB
			"detail_node": &graphql.Field{
				Type:        graphqlType.NodeDetailType,
				Description: "Get a node_detail by uuid",
				Args: graphql.FieldConfigArgument{
					"node_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node_detail")
					return dao.ReadNodeDetail(params.Args)
				},
			},
		},
	})
