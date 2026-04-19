LOCAL_BIN=$(CURDIR)/bin

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/protobuf/proto
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc


generate-api:
	mkdir -p pkg/user/v1
	protoc --proto_path api/user/v1 \
	--go_out=pkg/user/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user/v1/user.proto

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/server/main.go

