proto:
	protoc -I internal/user/proto --go_out=./internal/user --go-grpc_out=require_unimplemented_servers=false:./internal/user internal/user/proto/*.proto
	protoc -I internal/event/proto --go_out=./internal/event --go-grpc_out=require_unimplemented_servers=false:./internal/event internal/event/proto/*.proto