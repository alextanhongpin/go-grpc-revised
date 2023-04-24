.PHONY: proto server client

PROTO_PATH=proto
proto:
	@protoc \
		--proto_path=${PROTO_PATH} \
		--go_out=${PROTO_PATH} \
		--go_opt=paths=source_relative \
		--go-grpc_out=${PROTO_PATH} \
		--go-grpc_opt=paths=source_relative ${PROTO_PATH}/*.proto


init:
	@-brew install protobuf
	@-brew upgrade protobuf
	@go get -u google.golang.org/grpc
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc


server: # Run the server.
	@go run server/main.go


client: # Run the client.
	@go run client/main.go -name=$(name)


is_valid_protoc_version:
ifneq ($(shell protoc --version),libprotoc 3.14.0)
	$(error Expected version libprotoc 3.14.0, got $(shell protoc --version))
endif

check_protoc_version: is_valid_protoc_version
