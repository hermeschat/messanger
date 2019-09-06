proto:
	curl https://raw.githubusercontent.com/hermeschat/proto/master/api.proto > api.proto
	protoc --go_out=plugins=grpc:. api.proto && cp api.pb.go ./api
	rm -rf api.pb.go api.proto