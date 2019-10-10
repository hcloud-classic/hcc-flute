package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/flute/logger"
	"hcc/flute/mysql"
	"hcc/flute/types"
	"time"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			////////////////////////////// Node ///////////////////////////////
			/* Get (read) single node by uuid
			   http://192.168.110.240:7000/graphql?query={node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c"){uuid,bmc_mac_addr,bmc_ip,pxe_mac_addr,status,cpu_cores,memory,desc,created_at}}
			*/
			"node": &graphql.Field{
				Type:        nodeType,
				Description: "Get a node by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node")

					requestedUUID, ok := p.Args["uuid"].(string)
					if ok {
						node := new(types.Node)

						var uuid string
						var BMCmacAddr string
						var bmcIP string
						var pxeMACaddr string
						var status string
						var cpuCores int
						var memory int
						var desc string
						var createdAt time.Time

						sql := "select uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, `desc`, created_at from node where uuid = ?"
						err := mysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid, &BMCmacAddr, &bmcIP, &pxeMACaddr, &status, &cpuCores, &memory, &desc, &createdAt)
						if err != nil {
							logger.Logger.Println(err)
							return nil, nil
						}

						node.UUID = uuid
						node.BmcMacAddr = BMCmacAddr
						node.BmcIP = bmcIP
						node.PXEMacAddr = pxeMACaddr
						node.Status = status
						node.CPUCores = cpuCores
						node.Memory = memory
						node.Desc = desc
						node.CreatedAt = createdAt

						return node, nil
					}
					return nil, nil
				},
			},

			/* Get (read) node list
			   http://192.168.110.240:7000/graphql?query={list_node{uuid,bmc_mac_addr,bmc_ip,pxe_mac_addr,status,cpu_cores,memory,desc,created_at}}
			*/
			"list_node": &graphql.Field{
				Type:        graphql.NewList(nodeType),
				Description: "Get node list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_node")

					var nodes []types.Node

					var uuid string
					var BMCmacAddr string
					var bmcIP string
					var pxeMACaddr string
					var status string
					var cpuCores int
					var memory int
					var desc string
					var createdAt time.Time

					sql := "select uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, `desc`, created_at from node where active = 1"
					stmt, err := mysql.Db.Query(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()

					for stmt.Next() {
						err := stmt.Scan(&uuid, &BMCmacAddr, &bmcIP, &pxeMACaddr, &status, &cpuCores, &memory, &desc, &createdAt)
						if err != nil {
							logger.Logger.Println(err)
						}

						node := types.Node{UUID: uuid, BmcMacAddr: BMCmacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMACaddr, Status: status, CPUCores: cpuCores, Memory: memory, Desc: desc, CreatedAt: createdAt}

						logger.Logger.Println(node)
						nodes = append(nodes, node)
					}

					return nodes, nil
				},
			},

			////////////////////////////// Node Detail ///////////////////////////////
			/* Get (read) detail of a node by uuid
			   http://192.168.110.240:7000/graphql?query={node_detail(node_uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c"){node_uuid,cpu_model,cpu_processors,cpu_threads}}
			*/
			"node_detail": &graphql.Field{
				Type:        nodedetailType,
				Description: "Get detail of a node by uuid",
				Args: graphql.FieldConfigArgument{
					"node_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node_detail")

					requestedNodeUUID, ok := p.Args["node_uuid"].(string)
					if ok {
						nodeDetail := new(types.NodeDetail)

						var nodeUUID string
						var cpuModel string
						var cpuProcessors int
						var cpuThreads int

						sql := "select * from node_detail where node_uuid = ?"
						err := mysql.Db.QueryRow(sql, requestedNodeUUID).Scan(&nodeUUID, &cpuModel, &cpuProcessors, &cpuThreads)
						if err != nil {
							logger.Logger.Println(err)
							return nil, nil
						}

						nodeDetail.NodeUUID = nodeUUID
						nodeDetail.CPUModel = cpuModel
						nodeDetail.CPUProcessors = cpuProcessors
						nodeDetail.CPUThreads = cpuThreads

						return nodeDetail, nil
					}
					return nil, nil
				},
			},

			/* Get (read) node detail list
			   http://192.168.110.240:7000/graphql?query={list_node_detail{node_uuid,cpu_model,cpu_processors,cpu_threads}}
			*/
			"list_node_detail": &graphql.Field{
				Type:        graphql.NewList(nodedetailType),
				Description: "Get node detail list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_node_detail")

					var nodeDetails []types.NodeDetail
					var nodeUUID string
					var cpuModel string
					var cpuProcessors int
					var cpuThreads int

					sql := "select nd.* from node n, node_detail nd where n.uuid = nd.node_uuid and n.active = 1"
					stmt, err := mysql.Db.Query(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()

					for stmt.Next() {
						err := stmt.Scan(&nodeUUID, &cpuModel, &cpuProcessors, &cpuThreads)
						if err != nil {
							logger.Logger.Println(err)
						}

						nodeDetail := types.NodeDetail{NodeUUID: nodeUUID, CPUModel: cpuModel, CPUProcessors: cpuProcessors, CPUThreads: cpuThreads}

						logger.Logger.Println(nodeDetail)
						nodeDetails = append(nodeDetails, nodeDetail)
					}

					return nodeDetails, nil
				},
			},
		},
	})
