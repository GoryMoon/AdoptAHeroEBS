### Adopt A Hero Connection EBS



## Requirements
* [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation)
* `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26`
* `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1`

## Setup
To generate ProtoBuff and gRPC code run the following:
```shell
protoc --proto_path=protos --go_out=internal/protos --go_opt=paths=source_relative \
    --go-grpc_out=internal/protos --go-grpc_opt=paths=source_relative \
    protos/*.proto
```

## Run

You need to specify the following environment variables:
* SECRET - Base64 encoded data

The following environment variables are optional:
* HOST - ""
* PORT - 50051

Build
```shell
go build -o bin/server cmd/server/main.go
```

Start the server
```shell
SECRET=data ./bin/server 
```

# License
MIT
