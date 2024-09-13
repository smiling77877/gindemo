package main

import (
	"gindemo/webook/pkg/wego"
	"gindemo/webook/tag/grpc"
	"gindemo/webook/tag/ioc"
	"gindemo/webook/tag/repository/cache"
	"gindemo/webook/tag/repository/dao"
	"gindemo/webook/tag/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis, ioc.InitLogger, ioc.InitDB)

func Init() *wego.App {
	wire.Build(thirdProvider, cache.NewRedisTagCache, dao.NewGORMTagDAO,
		ioc.InitRepository, service.NewTagService, grpc.NewTagServiceServer,
		ioc.InitGRPCxServer, ioc.InitKafka, ioc.InitProducer, ioc.InitEtcdClient,
		wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
