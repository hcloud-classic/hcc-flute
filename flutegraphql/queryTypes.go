package flutegraphql

import (
	"GraphQL_Flute/flutelogger"
	"GraphQL_Flute/flutemysql"
	"GraphQL_Flute/flutetypes"
	"github.com/graphql-go/graphql"
	"time"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			////////////////////////////// Node ///////////////////////////////
			/* Get (read) single node by uuid
			   http://localhost:8001/graphql?query={node(uuid:"[node_uuid]"){uuid,mac_addr,ipmi_ip,status,cpu,memory,detail,created_at}}
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
					flutelogger.Logger.Println("Resolving: node")

					requestedUUID, ok := p.Args["uuid"].(string)
					if ok {
						node := new(flutetypes.Node)

						var uuid string
						var macAddr string
						var ipmiIp string
						var status string
						var cpu int
						var memory int
						var detail string
						var createdAt time.Time

						sql := "select * from node where uuid = ?"
						err := flutemysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid, &macAddr, &ipmiIp, &status, &cpu, &memory, &detail, &createdAt)
						if err != nil {
							flutelogger.Logger.Println(err)
							return nil, nil
						}

						node.UUID = uuid
						node.MacAddr = macAddr
						node.IpmiIP = ipmiIp
						node.Status = status
						node.Cpu = cpu
						node.Memory = memory
						node.Detail = detail
						node.CreatedAt = createdAt

						return node, nil
					}
					return nil, nil
				},
			},

			/* Get (read) node list
			   http://localhost:8001/graphql?query={list_node{uuid,mac_addr,ipmi_ip,status,cpu,memory,detail,created_at}}
			*/
			"list_node": &graphql.Field{
				Type:        graphql.NewList(nodeType),
				Description: "Get node list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					flutelogger.Logger.Println("Resolving: list_node")

					var nodes []flutetypes.Node
					var uuid string
					var macAddr string
					var ipmiIp string
					var status string
					var cpu int
					var memory int
					var detail string
					var createdAt time.Time

					sql := "select * from node"
					stmt, err := flutemysql.Db.Query(sql)
					if err != nil {
						flutelogger.Logger.Println(err)
						return nil, nil
					}
					defer stmt.Close()

					for stmt.Next() {
						err := stmt.Scan(&uuid, &macAddr, &ipmiIp, &status, &cpu, &memory, &detail, &createdAt)
						if err != nil {
							flutelogger.Logger.Println(err)
						}

						node := flutetypes.Node{UUID: uuid, MacAddr: macAddr, IpmiIP: ipmiIp, Status: status, Cpu: cpu, Memory: memory, Detail: detail, CreatedAt: createdAt}

						flutelogger.Logger.Println(node)
						nodes = append(nodes, node)
					}

					return nodes, nil
				},
			},

			////////////////////////////// Ipmi ///////////////////////////////
			/* Get (read) single ipmi by uuid
			   http://localhost:8001/graphql?query={ipmi(uuid:"[ipmi_uuid]]"){uuid}}
			*/
			"ipmi": &graphql.Field{
				Type:        nodeType,
				Description: "Get a ipmi by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					flutelogger.Logger.Println("Resolving: ipmi")

					requestedUUID, ok := p.Args["uuid"].(string)
					if ok {
						ipmi := new(flutetypes.Ipmi)

						var uuid string

						sql := "select * from ipmi where uuid = ?"
						err := flutemysql.Db.QueryRow(sql, requestedUUID).Scan(&uuid)
						if err != nil {
							flutelogger.Logger.Println(err)
							return nil, nil
						}

						ipmi.UUID = uuid

						return ipmi, nil
					}
					return nil, nil
				},
			},

			/* Get (read) ipmi list
			   http://localhost:8001/graphql?query={list_ipmi{uuid}}
			*/
			"list_ipmi": &graphql.Field{
				Type:        graphql.NewList(nodeType),
				Description: "Get ipmi list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					flutelogger.Logger.Println("Resolving: list_ipmi")

					var ipmis []flutetypes.Ipmi
					var uuid string

					sql := "select * from ipmi"
					stmt, err := flutemysql.Db.Query(sql)
					if err != nil {
						flutelogger.Logger.Println(err)
						return nil, nil
					}
					defer stmt.Close()

					for stmt.Next() {
						err := stmt.Scan(&uuid)
						if err != nil {
							flutelogger.Logger.Println(err)
						}

						ipmi := flutetypes.Ipmi{UUID: uuid}

						flutelogger.Logger.Println(ipmi)
						ipmis = append(ipmis, ipmi)
					}

					return ipmis, nil
				},
			},
		},
	})
