package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/flute/config"
	"hcc/flute/ipmi"
	"hcc/flute/logger"
	"hcc/flute/mysql"
	"hcc/flute/types"
	"strconv"
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

				bmcIP, bmcIPOk := params.Args["bmc_ip"].(string)

				if bmcIPOk {
					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					uuid, err := ipmi.GetUUID(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					BMCmac, err := ipmi.GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoBMC), true)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					PXEmac, err := ipmi.GetNICMac(bmcIP, int(config.Ipmi.BaseboardNICNoPXE), false)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					powerState, err := ipmi.GetPowerState(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					processors, err := ipmi.GetProcessors(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					cpuCores, err := ipmi.GetProcessorsCores(bmcIP, serialNo, processors)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					memory, err := ipmi.GetTotalSystemMemory(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					node := types.Node{
						UUID:        uuid,
						BmcMacAddr:  BMCmac,
						BmcIP:       bmcIP,
						PXEMacAddr:  PXEmac,
						Status:      powerState,
						CPUCores:    cpuCores,
						Memory:      memory,
						Description: params.Args["description"].(string),
					}

					sql := "insert into node(uuid, bmc_mac_addr, bmc_ip, pxe_mac_addr, status, cpu_cores, memory, description, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, now())"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(node.UUID, node.BmcMacAddr, node.BmcIP, node.PXEMacAddr, node.Status, node.CPUCores, node.Memory, node.Description)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, nil
					}
					logger.Logger.Println(result.LastInsertId())

					return node, nil
				}

				return nil, nil
			},
		},
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
						return nil, nil
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					state, _ := ipmi.GetPowerState(bmcIP, serialNo)
					if state == "On" {
						return "Already turned on", nil
					}

					result, err := ipmi.ChangePowerState(bmcIP, serialNo, "On")
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					return result, nil
				}

				return nil, nil
			},
		},
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
						return nil, nil
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
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
						return nil, nil
					}

					return result, nil
				}

				return nil, nil
			},
		},
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

				nodeUUID, nodeUUIDOk := params.Args["node_uuid"].(string)

				if nodeUUIDOk {
					var bmcIP string

					sql := "select bmc_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, nodeUUID).Scan(&bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					serialNo, err := ipmi.GetSerialNo(bmcIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					model, err := ipmi.GetProcessorModel(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					processors, err := ipmi.GetProcessors(bmcIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					threads, err := ipmi.GetProcessorsThreads(bmcIP, serialNo, processors)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
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
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(nodedetail.NodeUUID, nodedetail.CPUModel, nodedetail.CPUProcessors, nodedetail.CPUThreads)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, nil
					}
					logger.Logger.Println(result.LastInsertId())

					return nodedetail, nil
				}

				return nil, nil
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

				requestUUIDD, requestUUIDDOK := params.Args["uuid"].(string)
				bmcMacAddr, bmcMacAddrOk := params.Args["bmc_mac_addr"].(string)
				bmcIp, bmcIpOk := params.Args["bmc_ip"].(string)
				pxeMacAdr, pxeMacAdrOk := params.Args["pxe_mac_addr"].(string)
				status, statusOk := params.Args["status"].(string)
				cpuCores, cpuCoresOk := params.Args["cpu_cores"].(int)
				memory, memoryOk := params.Args["memory"].(int)
				description, descriptionOk := params.Args["description"].(string)
				active, activeOk := params.Args["active"].(int)

				node := new(types.Node)
				node.UUID = requestUUIDD
				node.BmcMacAddr = bmcMacAddr
				node.BmcIP = bmcIp
				node.PXEMacAddr = pxeMacAdr
				node.Status = status
				node.CPUCores = cpuCores
				node.Memory = memory
				node.Description = description
				node.Active = active

				if requestUUIDDOK {
					if !bmcMacAddrOk && !bmcIpOk && !pxeMacAdrOk && !statusOk && !cpuCoresOk && !memoryOk && !descriptionOk && !activeOk {
						return nil, nil
					}

					sql := "update node set"
					if bmcMacAddrOk {
						sql += " bmc_mac_addr = '" + bmcMacAddr + "'"
						if bmcIpOk || pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += ", "
						}
					}
					if bmcIpOk {
						sql += " bmc_ip = '" + bmcIp + "'"
						if pxeMacAdrOk || statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += ", "
						}
					}
					if pxeMacAdrOk {
						sql += " pxe_mac_addr = '" + pxeMacAdr + "'"
						if statusOk || cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += ", "
						}
					}
					if statusOk {
						sql += " status = '" + status + "'"
						if cpuCoresOk || memoryOk || descriptionOk || activeOk {
							sql += ", "
						}
					}
					if cpuCoresOk {
						sql += " cpu_cores = '" + strconv.Itoa(cpuCores) + "'"
						if memoryOk || descriptionOk || activeOk {
							sql += ", "
						}
					}
					if memoryOk {
						sql += " memory = '" + strconv.Itoa(memory) + "'"
						if descriptionOk || activeOk {
							sql += ", "
						}
					}
					if descriptionOk {
						sql += " description = '" + description + "'"
						if activeOk {
							sql += ", "
						}
					}
					if activeOk {
						sql += " active = '" + strconv.Itoa(active) + "'"
					}
					sql += " where uuid = ?"

					logger.Logger.Println("update_node sql : ", sql)
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer stmt.Close()

					result, err2 := stmt.Exec(node.UUID)
					if err2 != nil {
						logger.Logger.Println(err2)
						return nil, nil
					}
					logger.Logger.Println(result.LastInsertId())
					return node, nil
				}
				return nil, nil
			},
		},
	},
})
