package flutegraphql

import (
	"GraphQL_Flute/fluteipmi"
	"GraphQL_Flute/flutelogger"
	"GraphQL_Flute/flutemysql"
	"GraphQL_Flute/flutetypes"
	"github.com/graphql-go/graphql"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		////////////////////////////// node ///////////////////////////////
		/* Create new node
		http://localhost:8001/graphql?query=mutation+_{create_node(ipmi_ip:"172.31.0.1",detail:"Compute1"){uuid,mac_addr,ipmi_ip,status,cpu,memory,detail,created_at}}
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
				flutelogger.Logger.Println("Resolving: create_node")

				ipmiIp := params.Args["ipmi_ip"].(string)

				serialNo, err := fluteipmi.GetSerialNo(ipmiIp)
				if err != nil {
					flutelogger.Logger.Fatal(err)
				}

				uuid, err := fluteipmi.GetUuid(ipmiIp, serialNo)
				if err != nil {
					flutelogger.Logger.Fatal(err)
					return nil, nil
				}

				mac, err := fluteipmi.GetBMCNICMac(ipmiIp)
				if err != nil {
					flutelogger.Logger.Fatal(err)
					return nil, nil
				}

				powerState, err := fluteipmi.GetPowerState(ipmiIp, serialNo)
				if err != nil {
					flutelogger.Logger.Fatal(err)
					return nil, nil
				}

				processors, err := fluteipmi.GetProcessors(ipmiIp, serialNo)
				if err != nil {
					flutelogger.Logger.Fatal(err)
					return nil, nil
				}

				cores, err := fluteipmi.GetProcessorsCores(ipmiIp, serialNo, processors)
				if err != nil {
					flutelogger.Logger.Fatal(err)
					return nil, nil
				}

				memory, err := fluteipmi.GetTotalSystemMemory(ipmiIp, serialNo)
				if err != nil {
					flutelogger.Logger.Fatal(err)
					return nil, nil
				}

				////////////////////////////////////////////////////////
				// Get node info from RestfulAPI by IPMI
				node := flutetypes.Node {
					UUID:         uuid,
					MacAddr:      mac,
					IpmiIP:       ipmiIp,
					Status:       powerState,
					Cpu:          cores,
					Memory:       memory,
					Detail:       params.Args["detail"].(string),
				}

				sql := "insert into node(uuid, mac_addr, ipmi_ip, status, cpu, memory, detail, created_at) values (?, ?, ?, ?, ?, ?, ?, now())"
				stmt, err := flutemysql.Db.Prepare(sql)
				if err != nil {
					flutelogger.Logger.Println(err.Error())
					return nil, nil
				}
				defer stmt.Close()
				result, err2 := stmt.Exec(node.UUID, node.MacAddr, node.IpmiIP, node.Status, node.Cpu, node.Memory, node.Detail)
				if err2 != nil {
					flutelogger.Logger.Println(err2)
					return nil, nil
				}
				flutelogger.Logger.Println(result.LastInsertId())

				return node, nil
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
		//		flutelogger.Logger.Println("Resolving: update_volume")
		//
		//		requestedUUID, _ := params.Args["uuid"].(string)
		//		size, sizeOk := params.Args["size"].(int)
		//		_type, _typeOk := params.Args["type"].(string)
		//		serverUUID, serverUUIDOk := params.Args["server_uuid"].(string)
		//
		//		volume := new(flutetypes.Volume)
		//
		//		if sizeOk && _typeOk && serverUUIDOk {
		//			volume.UUID = requestedUUID
		//			volume.Size = size
		//			volume.Type = _type
		//			volume.ServerUUID = serverUUID
		//
		//			sql := "update volume set size = ?, type = ?, server_uuid = ? where uuid = ?"
		//			stmt, err := flutemysql.Db.Prepare(sql)
		//			if err != nil {
		//				flutelogger.Logger.Println(err.Error())
		//				return nil, nil
		//			}
		//			defer stmt.Close()
		//			result, err2 := stmt.Exec(volume.Size, volume.Type, volume.ServerUUID, volume.UUID)
		//			if err2 != nil {
		//				flutelogger.Logger.Println(err2)
		//				return nil, nil
		//			}
		//			flutelogger.Logger.Println(result.LastInsertId())
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
		//		flutelogger.Logger.Println("Resolving: delete_volume")
		//
		//		requestedUUID, ok := params.Args["uuid"].(string)
		//		if ok {
		//			sql := "delete from volume where uuid = ?"
		//			stmt, err := flutemysql.Db.Prepare(sql)
		//			if err != nil {
		//				flutelogger.Logger.Println(err.Error())
		//				return nil, nil
		//			}
		//			defer stmt.Close()
		//			result, err2 := stmt.Exec(requestedUUID)
		//			if err2 != nil {
		//				flutelogger.Logger.Println(err2)
		//				return nil, nil
		//			}
		//			flutelogger.Logger.Println(result.RowsAffected())
		//
		//			return requestedUUID, nil
		//		}
		//		return nil, nil
		//	},
		//},
	},
})
