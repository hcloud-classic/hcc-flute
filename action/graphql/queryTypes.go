package graphql

import (
<<<<<<< HEAD
	"github.com/graphql-go/graphql"
	graphqlType "hcc/flute/action/graphql/type"
	"hcc/flute/dao"
	"hcc/flute/lib/logger"
=======
	"errors"
	"hcc/flute/lib/logger"
	"hcc/flute/lib/mysql"
	"hcc/flute/model"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
>>>>>>> f41ff24 (Refactoring packages structure)
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
<<<<<<< HEAD
			// IPMI
			"power_state_node": &graphql.Field{
				Type:        graphql.String,
				Description: "Get the node's power state",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: power_state_node")
					return dao.GetPowerStateNode(params.Args)
				},
			},
			// node DB
			"node": &graphql.Field{
				Type:        graphqlType.NodeType,
=======
			"node": &graphql.Field{
				Type:        nodeType,
>>>>>>> f41ff24 (Refactoring packages structure)
				Description: "Get a node by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
<<<<<<< HEAD
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
=======
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node")

					requestedUUID, ok := p.Args["uuid"].(string)
					if ok {
						node := new(model.Node)

						var uuid string
						var BMCmacAddr string
						var bmcIP string
						var pxeMacAddr string
						var status string
						var cpuCores int
						var memory int
						var description string
						var createdAt time.Time
						var active int

						sql := "select * from node where uuid = ?"
						err := mysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid, &BMCmacAddr, &bmcIP, &pxeMacAddr, &status, &cpuCores, &memory, &description, &createdAt, &active)
						if err != nil {
							logger.Logger.Println(err)
							return nil, err
						}

						node.UUID = uuid
						node.BmcMacAddr = BMCmacAddr
						node.BmcIP = bmcIP
						node.PXEMacAddr = pxeMacAddr
						node.Status = status
						node.CPUCores = cpuCores
						node.Memory = memory
						node.Description = description
						node.CreatedAt = createdAt
						node.Active = active

						return node, nil
					}
					return nil, errors.New("need uuid argument")
				},
			},
			"list_node": &graphql.Field{
				Type:        graphql.NewList(nodeType),
				Description: "Get node list",
				Args: graphql.FieldConfigArgument{
>>>>>>> f41ff24 (Refactoring packages structure)
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
<<<<<<< HEAD
						Type: graphql.Int,
=======
						Type: graphql.String,
>>>>>>> f41ff24 (Refactoring packages structure)
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
<<<<<<< HEAD
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
=======

					var nodes []model.Node
					var uuid string
					var createdAt time.Time

					bmcMacAddr, bmcMacAddrOk := params.Args["bmc_mac_addr"].(string)
					bmcIP, bmcIPOk := params.Args["bmc_ip"].(string)
					pxeMacAdr, pxeMacAdrOk := params.Args["pxe_mac_addr"].(string)
					status, statusOk := params.Args["status"].(string)
					cpuCores, cpuCoresOk := params.Args["cpu_cores"].(int)
					memory, memoryOk := params.Args["memory"].(int)
					description, descriptionOk := params.Args["description"].(string)
					active, activeOk := params.Args["active"].(int)
					row, rowOk := params.Args["row"].(int)
					page, pageOk := params.Args["page"].(int)
					if !rowOk || !pageOk {
						return nil, nil
					}

					sql := "select * from node where"
					if bmcMacAddrOk {
						sql += " bmc_mac_addr = '" + bmcMacAddr + "'"
						if bmcIPOk || pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += " and"
						}
					}
					if bmcIPOk {
						sql += " bmc_ip = '" + bmcIP + "'"
						if pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += " and"
						}
					}
					if pxeMacAdrOk {
						sql += " pxe_mac_addr = '" + pxeMacAdr + "'"
						if statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += " and"
						}
					}
					if statusOk {
						sql += " status = '" + status + "'"
						if cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += " and"
						}
					}
					if cpuCoresOk {
						sql += " cpu_cores = '" + strconv.Itoa(cpuCores) + "'"
						if memoryOk || descriptionOk || activeOk {
							sql += " and"
						}
					}
					if memoryOk {
						sql += " memory = '" + strconv.Itoa(memory) + "'"
						if descriptionOk || activeOk {
							sql += " and"
						}
					}
					if descriptionOk {
						sql += " description = '" + description + "'"
						if activeOk {
							sql += " and"
						}
					}
					if activeOk {
						sql += " active = '" + strconv.Itoa(active) + "'"
					}
					sql += " order by created_at desc limit ? offset ?"

					logger.Logger.Println("list_node sql : ", sql)

					stmt, err := mysql.Db.Query(sql, row, row*(page-1))
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()

					for stmt.Next() {
						err := stmt.Scan(&uuid, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &createdAt, &active)
						if err != nil {
							logger.Logger.Println(err)
						}
						node := model.Node{UUID: uuid, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, CreatedAt: createdAt, Active: active}
						nodes = append(nodes, node)
					}
					return nodes, nil
				},
			},
			"all_node": &graphql.Field{
				Type:        graphql.NewList(nodeType),
				Description: "Get all node list",
				Args: graphql.FieldConfigArgument{
>>>>>>> f41ff24 (Refactoring packages structure)
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_node")
<<<<<<< HEAD
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
=======

					var nodes []model.Node
					var uuid string
					var bmcMacAddr string
					var bmcIP string
					var pxeMacAdr string
					var status string
					var cpuCores int
					var memory int
					var description string
					var createdAt time.Time
					var active int
					row, rowOk := params.Args["row"].(int)
					page, pageOk := params.Args["page"].(int)
					if !rowOk || !pageOk {
						return nil, nil
					}

					sql := "select * from node order by created_at desc limit ? offset ?"
					logger.Logger.Println("list_server sql  : ", sql)
					stmt, err := mysql.Db.Query(sql, row, row*(page-1))
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()

					for stmt.Next() {
						err := stmt.Scan(&uuid, &bmcMacAddr, &bmcIP, &pxeMacAdr, &status, &cpuCores, &memory, &description, &createdAt, &active)
						if err != nil {
							logger.Logger.Println(err)
							return nil, err
						}
						node := model.Node{UUID: uuid, BmcMacAddr: bmcMacAddr, BmcIP: bmcIP, PXEMacAddr: pxeMacAdr, Status: status, CPUCores: cpuCores, Memory: memory, Description: description, CreatedAt: createdAt, Active: active}
						nodes = append(nodes, node)
					}
					return nodes, nil
				},
			},
			"detail_node": &graphql.Field{
				Type:        nodeDetailType,
				Description: "Get detail of a node by uuid",
>>>>>>> f41ff24 (Refactoring packages structure)
				Args: graphql.FieldConfigArgument{
					"node_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
<<<<<<< HEAD
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node_detail")
					return dao.ReadNodeDetail(params.Args)
=======
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: node_detail")

					nodeDetail := new(model.NodeDetail)
					var nodeUUID string
					var cpuModel string
					var cpuProcessors int
					var cpuThreads int
					requestedNodeUUID, requestedNodeUUIDok := p.Args["node_uuid"].(string)
					if requestedNodeUUIDok {
						sql := "select * from node_detail where node_uuid = ?"
						err := mysql.Db.QueryRow(sql, requestedNodeUUID).Scan(&nodeUUID, &cpuModel, &cpuProcessors, &cpuThreads)
						if err != nil {
							logger.Logger.Println(err)
							return nil, err
						}
						nodeDetail.NodeUUID = nodeUUID
						nodeDetail.CPUModel = cpuModel
						nodeDetail.CPUProcessors = cpuProcessors
						nodeDetail.CPUThreads = cpuThreads

						return nodeDetail, nil
					}
					return nil, errors.New("need node_uuid argument")
				},
			},
			"num_node": &graphql.Field{
				Type:        nodeNum,
				Description: "Get the number of node",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: num_node")

					var nodeNum model.NodeNum
					var nodeNr int

					sql := "select count(*) from node"
					err := mysql.Db.QueryRow(sql).Scan(&nodeNr)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					logger.Logger.Println("Count: ", nodeNr)
					nodeNum.Number = nodeNr

					return nodeNum, nil
>>>>>>> f41ff24 (Refactoring packages structure)
				},
			},
		},
	})
