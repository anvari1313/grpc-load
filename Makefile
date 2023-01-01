generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/proto.proto

mod:
	go mod download

build-linux-http:
	go build .
