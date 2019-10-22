package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"hcc/flute/dao"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/wol"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create_node": &graphql.Field{
			Type:        nodeType,
			Description: "Create new node",
			Args: graphql.FieldConfigArgument{
				"bmc_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_node")

				return dao.CreateNode(params.Args)
			},
		},
		//"on_node": &graphql.Field{
		//	Type:        graphql.String,
		//	Description: "On node",
		//	Args: graphql.FieldConfigArgument{
		//		"uuid": &graphql.ArgumentConfig{
		//			Type: graphql.String,
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		logger.Logger.Println("Resolving: on_node")
		//
		//		return dao.OnNode(params.Args)
		//	},
		//},
		"on_node": &graphql.Field{
			Type:        graphql.String,
			Description: "On node",
			Args: graphql.FieldConfigArgument{
				"mac": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: on_node")
				mac, macOk := params.Args["mac"].(string)
				if macOk {
					err := wol.OnNode(mac)
					if err != nil {
						return nil, err
					}

					return "Send magic packet to \"" + mac + "\"", nil
				}

				return nil, errors.New("need mac argument")
			},
		},
		//"off_node": &graphql.Field{
		//	Type:        graphql.String,
		//	Description: "Off node",
		//	Args: graphql.FieldConfigArgument{
		//		"uuid": &graphql.ArgumentConfig{
		//			Type: graphql.String,
		//		},
		//		"force_off": &graphql.ArgumentConfig{
		//			Type: graphql.Boolean,
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		logger.Logger.Println("Resolving: off_node")
		//
		//		return dao.OffNode(params.Args)
		//	},
		//},
		"create_node_detail": &graphql.Field{
			Type:        nodeDetailType,
			Description: "Create new node_detail",
			Args: graphql.FieldConfigArgument{
				"node_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_node_uuid")

				return dao.CreateNodeDetail(params.Args)
			},
		},
		"update_node": &graphql.Field{
			Type:        nodeType,
			Description: "Update node",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
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
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_node")

				return dao.UpdateNode(params.Args)
			},
		},
	},
})
