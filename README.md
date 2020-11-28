# go-grpc-revised

Revised example of gRPC with golang.


## Quickstart

```bash
# Install all the required dependencies.
$ make init

$ protoc --version
libprotoc 3.14.0

$ make check_protoc_version

# Regenerate the go code from the proto file.
$ make proto

# Start the gRPC server.
$ make server

# Start the gRPC client that calls the server.
$ make client
```
