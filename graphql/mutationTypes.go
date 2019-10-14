package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"hcc/flute/config"
	"hcc/flute/ipmi"
	"hcc/flute/logger"
	"hcc/flute/mysql"
	"hcc/flute/types"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		////////////////////////////// node ///////////////////////////////
		/* Create new node
		http://192.168.110.240:7000/graphql?query=mutation+_{create_node(bmc_ip:"172.31.0.1",desc:"Compute1"){uuid,bmc_mac_addr,bmc_ip,status,cpu_cores,memory,desc,created_at}}
		*/
		"create_node": &graphql.Field{
			Type:        nodeType,
			Description: "Create new node",
			Args: graphql.FieldConfigArgument{
				"bmc_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"desc": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_node")

				bmcIP, bmcIPOk := params.Args["bmc_ip"].(string)
				desc, descOk := params.Args["desc"].(string)

				if !descOk {
					desc = ""
				}

				if bmcIPOk {
					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					uuid, err := ipmi.GetUUID(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					BMCmac, err := ipmi.GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoBMC), true)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					PXEmac, err := ipmi.GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoPXE), false)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					powerState, err := ipmi.GetPowerState(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					processors, err := ipmi.GetProcessors(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					cpuCores, err := ipmi.GetProcessorsCores(bmcIP, serialNo, processors)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					memory, err := ipmi.GetTotalSystemMemory(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					node := types.Node{
						UUID:       uuid,
						BmcMacAddr: BMCmac,
						BmcIP:      bmcIP,
						PXEMacAddr: PXEmac,
						Status:     powerState,
						CPUCores:   cpuCores,
						Memory:     memory,
						Desc:       desc,
					}

					sql := "insert into node(uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, `desc`, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, now())"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(node.UUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Desc)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, err2
					}
					logger.Logger.Println(result.LastInsertId())

					err = ipmi.BMCIPParserCheckActive(node.BmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					return node, nil
				}

				return nil, errors.New("need bmc_ip argument")
			},
		},

		/* Update node desc
		http://192.168.110.240:7000/graphql?query=mutation+_{update_node_desc(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c",desc:"Compute2"){uuid,desc}}
		*/
		"update_node_desc": &graphql.Field{
			Type:        nodeType,
			Description: "Update node desc",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"desc": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_node_desc")

				uuid, uuidOk := params.Args["uuid"].(string)
				desc, descOk := params.Args["desc"].(string)

				if uuidOk && descOk {

					node := types.Node{
						UUID: uuid,
						Desc: desc,
					}

					sql := "update node set `desc` = ? where uuid = ?"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(node.Desc, node.UUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, err2
					}
					logger.Logger.Println(result.LastInsertId())

					return node, nil
				}

				return nil, errors.New("need uuid and desc arguments")
			},
		},

		/* Update all infos of all nodes (except power state)
		http://192.168.110.240:7000/graphql?query=mutation+_{update_all_nodes{uuid,bmc_mac_addr,bmc_ip,status,cpu_cores,memory,desc,created_at}}
		*/
		"update_all_nodes": &graphql.Field{
			Type:        graphql.NewList(nodeType),
			Description: "Update all infos of the node",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_all_nodes")

				return ipmi.UpdateAllNodes()
			},
		},

		/* Update status of the node
		http://192.168.110.240:7000/graphql?query=mutation+_{update_status_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c"){status}}
		*/
		"update_status_node": &graphql.Field{
			Type:        nodeType,
			Description: "Update status of the node",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_status_node")

				uuid, uuidOk := params.Args["uuid"].(string)

				if uuidOk {
					var bmcIP string

					sql := "select bmc_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					powerState, err := ipmi.GetPowerState(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					node := types.Node{
						UUID:   uuid,
						Status: powerState,
					}

					sql = "update node set status = ? where uuid = ?"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(node.Status, node.UUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, err2
					}
					logger.Logger.Println(result.LastInsertId())

					return node, nil
				}

				return nil, errors.New("need uuid argument")
			},
		},

		/* Update status of all nodes
		http://192.168.110.240:7000/graphql?query=mutation+_{update_status_nodes{status}}
		*/
		"update_status_nodes": &graphql.Field{
			Type:        graphql.NewList(nodeType),
			Description: "Update status of all nodes",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_status_nodes")

				return ipmi.UpdateStatusNodes()
			},
		},

		/* On node
		http://192.168.110.240:7000/graphql?query=mutation+_{on_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c")}
		*/
		"on_node": &graphql.Field{
			Type:        graphql.String,
			Description: "On node",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: on_node")

				uuid, uuidOk := params.Args["uuid"].(string)

				if uuidOk {
					var bmcIP string

					sql := "select bmc_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					state, _ := ipmi.GetPowerState(bmcIP, serialNo)
					if state == "On" {
						return "Already turned on", nil
					}

					result, err := ipmi.ChangePowerState(bmcIP, serialNo, "On")
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					return result, nil
				}

				return nil, errors.New("need uuid argument")
			},
		},

		/* Off node
		http://192.168.110.240:7000/graphql?query=mutation+_{off_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c")}
		http://192.168.110.240:7000/graphql?query=mutation+_{off_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c",force_off:true)}
		*/
		"off_node": &graphql.Field{
			Type:        graphql.String,
			Description: "Off node",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"force_off": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: off_node")

				uuid, uuidOk := params.Args["uuid"].(string)
				forceOff, _ := params.Args["force_off"].(bool)

				if uuidOk {
					var bmcIP string

					sql := "select bmc_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, uuid).Scan(&bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					state, _ := ipmi.GetPowerState(bmcIP, serialNo)
					if state == "Off" {
						return "Already turned off", nil
					}

					changeState := "GracefulShutdown"
					if forceOff {
						changeState = "ForceOff"
					}
					result, err := ipmi.ChangePowerState(bmcIP, serialNo, changeState)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					return result, nil
				}

				return nil, errors.New("need uuid argument")
			},
		},

		////////////////////////////// node_detail ///////////////////////////////
		/* Create new node_detail
		http://192.168.110.240:7000/graphql?query=mutation+_{create_node_detail(node_uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c"){node_uuid,cpu_model,cpu_processors,cpu_threads}}
		*/
		"create_node_detail": &graphql.Field{
			Type:        nodedetailType,
			Description: "Create new node_detail",
			Args: graphql.FieldConfigArgument{
				///////////////////////////////////////////
				"node_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_node_uuid")

				nodeUUID, nodeUUIDOk := params.Args["node_uuid"].(string)

				if nodeUUIDOk {
					var bmcIP string

					sql := "select bmc_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, nodeUUID).Scan(&bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					model, err := ipmi.GetProcessorModel(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					processors, err := ipmi.GetProcessors(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					threads, err := ipmi.GetProcessorsThreads(bmcIP, serialNo, processors)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}

					nodedetail := types.NodeDetail{
						NodeUUID:      nodeUUID,
						CPUModel:      model,
						CPUProcessors: processors,
						CPUThreads:    threads,
					}

					sql = "insert into node_detail(node_uuid, cpu_model, cpu_processors, cpu_threads) values (?, ?, ?, ?)"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err)
						return nil, err
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(nodedetail.NodeUUID, nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, err2
					}
					logger.Logger.Println(result.LastInsertId())

					return nodedetail, nil
				}

				return nil, errors.New("need node_uuid argument")
			},
		},

		/* Update detail infos of all nodes
		http://192.168.110.240:7000/graphql?query=mutation+_{update_nodes_detail(){status}}
		*/
		"update_nodes_detail": &graphql.Field{
			Type:        graphql.NewList(nodedetailType),
			Description: "Update detail infos of all nodes",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_nodes_detail")

				return ipmi.UpdateNodesDetail()
			},
		},
	},
})
