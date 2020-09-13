package monitoring

import (
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

//RegisterGRPCMetrics registers GRPC metrics in prometheus
func RegisterGRPCMetrics(server *grpc.Server) {
	grpc_prom.Register(server)
}