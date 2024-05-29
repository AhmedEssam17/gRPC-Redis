module grpc-redis/server

go 1.22.3

require (
	github.com/go-redis/redis v6.15.9+incompatible
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
	grpc-redis-server/protos/todo v0.0.0
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.33.1 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
)

replace grpc-redis-server/protos/todo => ../protos/todo
