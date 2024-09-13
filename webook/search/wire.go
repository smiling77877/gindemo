//go:build wireinject

package main

import (
	"gindemo/webook/search/events"
	"gindemo/webook/search/grpc"
	"gindemo/webook/search/ioc"
	"gindemo/webook/search/repository"
	"gindemo/webook/search/repository/dao"
	"gindemo/webook/search/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewUserElasticDAO, dao.NewArticleElasticDAO, dao.NewAnyESDAO, dao.NewTagESDAO,
	repository.NewUserRepository, repository.NewArticleRepository, repository.NewAnyRepository,
	service.NewSearchService, service.NewSyncService)

var thirdProvider = wire.NewSet(
	ioc.InitESClient, ioc.InitEtcdClient, ioc.InitLogger, ioc.InitKafka)

func Init() *App {
	wire.Build(thirdProvider, serviceProviderSet,
		grpc.NewSyncServiceServer, grpc.NewSearchService,
		events.NewUserConsumer, events.NewArticleConsumer,
		ioc.InitGRPCxServer, ioc.NewConsumers,
		wire.Struct(new(App), "*"))
	return new(App)
}
