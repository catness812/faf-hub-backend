proto:
	protoc -I internal/user/proto --go_out=./internal/user --go-grpc_out=require_unimplemented_servers=false:./internal/user internal/user/proto/*.proto
	protoc -I internal/event/proto --go_out=./internal/event --go-grpc_out=require_unimplemented_servers=false:./internal/event internal/event/proto/*.proto
	protoc -I internal/content/proto --go_out=./internal/content --go-grpc_out=require_unimplemented_servers=false:./internal/content internal/content/proto/*.proto
	protoc -I internal/notification/proto --go_out=./internal/notification --go-grpc_out=require_unimplemented_servers=false:./internal/notification internal/notification/proto/*.proto