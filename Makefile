proto:
	protoc --go_out=plugins=grpc:. ./api/pb/api.proto