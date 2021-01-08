module hcc/flute

go 1.14

require (
	github.com/Terry-Mao/goconf v0.0.0-20161115082538-13cb73d70c44
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/hcloud-classic/hcc_errors v1.1.0
	github.com/hcloud-classic/pb v0.0.0
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/hcloud-classic/pb => ../pb
