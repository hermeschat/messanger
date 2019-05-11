gen-proto:
	protoc --go_out=plugins=grpc:. ./pkg/api/*.proto