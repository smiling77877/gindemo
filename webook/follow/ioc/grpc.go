package ioc

import (
	grpc2 "gindemo/webook/follow/grpc"
	"gindemo/webook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(followRelation *grpc2.FollowServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	followRelation.Register(server)
	return &grpcx.Server{
		Server: server,
	}
}
