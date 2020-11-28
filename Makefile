.PHONY: proto server client

PROTO_PATH=proto
proto:
	@protoc --proto_path=${PROTO_PATH} --go_out=${PROTO_PATH} --go_opt=paths=source_relative --go-grpc_out=${PROTO_PATH} --go-grpc_opt=paths=source_relative ${PROTO_PATH}/*.proto


init:
	@brew upgrade protobuf
	@go get -u google.golang.org/grpc
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc


server:
	@go run server/main.go

client:
	@go run client/main.go
