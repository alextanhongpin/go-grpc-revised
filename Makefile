.PHONY: proto server client

PROTO_PATH=proto
proto:
	@find ${PROTO_PATH}/*.protoc -exec protoc -I ${PROTO_PATH} {} --go_out=plugins=grpc:${PROTO_PATH} \;


install:
	@brew install protobuf
	@brew upgrade protobuf
	@go get -u google.golang.org/grpc
	@go get -u github.com/golang/protobuf/protoc-gen-go


server:
	@go run server/main.go

client:
	@go run client/main.go
