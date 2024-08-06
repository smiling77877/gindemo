package ioc

import (
	grpc2 "gindemo/webook/interactive/grpc"
	"gindemo/webook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	addr := viper.GetString("grpc.server.addr")
	return &grpcx.Server{
		Server: s,
		Addr:   addr,
	}
}
