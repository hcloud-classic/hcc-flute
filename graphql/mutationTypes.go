package graphql

import (
	"github.com/graphql-go/graphql"
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
		http://localhost:7000/graphql?query=mutation+_{create_node(ipmi_ip:"172.31.0.1",detail:"Compute1"){uuid,mac_addr,ipmi_ip,status,cpu,memory,detail,created_at}}
		*/
		"create_node": &graphql.Field{
			Type:        nodeType,
			Description: "Create new node",
			Args: graphql.FieldConfigArgument{
				"ipmi_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"detail": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: create_node")

				ipmiIP, uuidOk := params.Args["ipmi_ip"].(string)

				if uuidOk {
					serialNo, err := ipmi.GetSerialNo(ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					uuid, err := ipmi.GetUUID(ipmiIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					mac, err := ipmi.GetBMCNICMac(ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					powerState, err := ipmi.GetPowerState(ipmiIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					processors, err := ipmi.GetProcessors(ipmiIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					cores, err := ipmi.GetProcessorsCores(ipmiIP, serialNo, processors)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					memory, err := ipmi.GetTotalSystemMemory(ipmiIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					node := types.Node{
						UUID:    uuid,
						MacAddr: mac,
						IpmiIP:  ipmiIP,
						Status:  powerState,
						CPU:     cores,
						Memory:  memory,
						Detail:  params.Args["detail"].(string),
					}

					sql := "insert into node(uuid, mac_addr, ipmi_ip, status, cpu, memory, detail, created_at) values (?, ?, ?, ?, ?, ?, ?, now())"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(node.UUID, node.MacAddr, node.IpmiIP, node.Status, node.CPU, node.Memory, node.Detail)
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
		//
		///* Update volume by uuid
		//   http://localhost:8001/graphql?query=mutation+_{update_volume(uuid:"[volume_uuid]",size:10240,type:"ext4",server_uuid:"[server_uuid]"){uuid,size,type,server_uuid}}
		//*/
		//"update_volume": &graphql.Field{
		//	Type:        volumeType,
		//	Description: "Update volume by uuid",
		//	Args: graphql.FieldConfigArgument{
		//		"uuid": &graphql.ArgumentConfig{
		//			Type: graphql.NewNonNull(graphql.String),
		//		},
		//		"size": &graphql.ArgumentConfig{
		//			Type: graphql.Int,
		//		},
		//		"type": &graphql.ArgumentConfig{
		//			Type: graphql.String,
		//		},
		//		"server_uuid": &graphql.ArgumentConfig{
		//			Type: graphql.String,
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		logger.Logger.Println("Resolving: update_volume")
		//
		//		requestedUUID, _ := params.Args["uuid"].(string)
		//		size, sizeOk := params.Args["size"].(int)
		//		_type, _typeOk := params.Args["type"].(string)
		//		serverUUID, serverUUIDOk := params.Args["server_uuid"].(string)
		//
		//		volume := new(types.Volume)
		//
		//		if sizeOk && _typeOk && serverUUIDOk {
		//			volume.UUID = requestedUUID
		//			volume.Size = size
		//			volume.Type = _type
		//			volume.ServerUUID = serverUUID
		//
		//			sql := "update volume set size = ?, type = ?, server_uuid = ? where uuid = ?"
		//			stmt, err := mysql.Db.Prepare(sql)
		//			if err != nil {
		//				logger.Logger.Println(err.Error())
		//				return nil, nil
		//			}
		//			defer stmt.Close()
		//			result, err2 := stmt.Exec(volume.Size, volume.Type, volume.ServerUUID, volume.UUID)
		//			if err2 != nil {
		//				logger.Logger.Println(err2)
		//				return nil, nil
		//			}
		//			logger.Logger.Println(result.LastInsertId())
		//
		//			return volume, nil
		//		}
		//		return nil, nil
		//	},
		//},
		//
		///* Delete volume by id
		//   http://localhost:8001/graphql?query=mutation+_{delete_volume(id:"test1"){id}}
		//*/
		//"delete_volume": &graphql.Field{
		//	Type:        volumeType,
		//	Description: "Delete volume by uuid",
		//	Args: graphql.FieldConfigArgument{
		//		"uuid": &graphql.ArgumentConfig{
		//			Type: graphql.NewNonNull(graphql.String),
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		logger.Logger.Println("Resolving: delete_volume")
		//
		//		requestedUUID, ok := params.Args["uuid"].(string)
		//		if ok {
		//			sql := "delete from volume where uuid = ?"
		//			stmt, err := mysql.Db.Prepare(sql)
		//			if err != nil {
		//				logger.Logger.Println(err.Error())
		//				return nil, nil
		//			}
		//			defer stmt.Close()
		//			result, err2 := stmt.Exec(requestedUUID)
		//			if err2 != nil {
		//				logger.Logger.Println(err2)
		//				return nil, nil
		//			}
		//			logger.Logger.Println(result.RowsAffected())
		//
		//			return requestedUUID, nil
		//		}
		//		return nil, nil
		//	},
		//},

		/* Update all infos of all nodes
		http://localhost:7000/graphql?query=mutation+_{update_all_nodes(){uuid,mac_addr,ipmi_ip,status,cpu,memory,detail,created_at}}
		*/
		"update_all_nodes": &graphql.Field{
			Type:        nodeType,
			Description: "Update all infos of the node",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_all_nodes")

				return ipmi.UpdateAllNodes()
			},
		},

		/* Update status of the node
		http://localhost:7000/graphql?query=mutation+_{update_status_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c"){status}}
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
					var ipmiIP string

					sql := "select ipmi_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, uuid).Scan(&ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					serialNo, err := ipmi.GetSerialNo(ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					powerState, err := ipmi.GetPowerState(ipmiIP, serialNo)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					node := types.Node{
						UUID:   uuid,
						Status: powerState,
					}

					sql = "update node set status = ? where uuid = ?"
					stmt, err := mysql.Db.Prepare(sql)
					if err != nil {
						logger.Logger.Println(err.Error())
						return nil, nil
					}
					defer func() {
						_ = stmt.Close()
					}()
					result, err2 := stmt.Exec(node.Status, node.UUID)
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

		/* Update status of all nodes
		http://localhost:7000/graphql?query=mutation+_{update_status_nodes(){status}}
		*/
		"update_status_nodes": &graphql.Field{
			Type:        nodeType,
			Description: "Update status of all nodes",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_status_nodes")

				return ipmi.UpdateStatusNodes()
			},
		},

		/* On node
		http://localhost:7000/graphql?query=mutation+_{on_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c")}
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
					var ipmiIP string

					sql := "select ipmi_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, uuid).Scan(&ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					serialNo, err := ipmi.GetSerialNo(ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					state, _ := ipmi.GetPowerState(ipmiIP, serialNo)
					if state == "On" {
						return "Already turned on", nil
					}

					result, err := ipmi.ChangePowerState(ipmiIP, serialNo, "On")
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					return result, nil
				}

				return nil, nil
			},
		},

		/* Off node
		http://localhost:7000/graphql?query=mutation+_{off_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c")}
		http://localhost:7000/graphql?query=mutation+_{off_node(uuid:"d4f3a900-b674-11e8-906e-000ffee02d5c",force_off:true)}
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
					var ipmiIP string

					sql := "select ipmi_ip from node where uuid = ?"
					err := mysql.Db.QueryRow(sql, uuid).Scan(&ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					serialNo, err := ipmi.GetSerialNo(ipmiIP)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					state, _ := ipmi.GetPowerState(ipmiIP, serialNo)
					if state == "Off" {
						return "Already turned off", nil
					}

					changeState := "GracefulShutdown"
					if forceOff {
						changeState = "ForceOff"
					}
					result, err := ipmi.ChangePowerState(ipmiIP, serialNo, changeState)
					if err != nil {
						logger.Logger.Println(err)
						return nil, nil
					}

					return result, nil
				}

				return nil, nil
			},
		},
	},
})
