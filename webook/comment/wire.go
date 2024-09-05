//go:build wireinject

package main

import (
	grpc2 "gindemo/webook/comment/grpc"
	"gindemo/webook/comment/ioc"
	"gindemo/webook/comment/repository"
	"gindemo/webook/comment/repository/dao"
	"gindemo/webook/comment/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO, repository.NewCommentRepo, service.NewCommentSvc, grpc2.NewGrpcServer)

var thirdProvider = wire.NewSet(ioc.InitLogger, ioc.InitDB)

func Init() *App {
	wire.Build(
		serviceProviderSet,
		thirdProvider,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"))
	return new(App)
}
