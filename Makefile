proto:
	protoc --go_out=plugins=grpc:. ./api/pb/api.proto && cp ./api/pb/api.pb.go ./example-client


	//https://raw.githubusercontent.com/hermeschat/proto/master/api.proto