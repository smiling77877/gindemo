//go:build wireinject

package main

import (
	"gindemo/webook/interactive/events"
	"gindemo/webook/interactive/grpc"
	"gindemo/webook/interactive/ioc"
	repository2 "gindemo/webook/interactive/repository"
	cache2 "gindemo/webook/interactive/repository/cache"
	dao2 "gindemo/webook/interactive/repository/dao"
	service2 "gindemo/webook/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
