//go:build wireinject

package main

import (
	"gindemo/webook/pkg/wego"
	"gindemo/webook/reward/grpc"
	"gindemo/webook/reward/ioc"
	"gindemo/webook/reward/repository"
	"gindemo/webook/reward/repository/cache"
	"gindemo/webook/reward/repository/dao"
	"gindemo/webook/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitEtcdClient, ioc.InitRedis)

func Init() *wego.App {
	wire.Build(thirdPartySet, service.NewWechatNativeRewardService,
		ioc.InitAccountClient, ioc.InitGRPCxServer, ioc.InitPaymentClient,
		repository.NewRewardRepository, cache.NewRewardRedisCache,
		dao.NewRewardGORMDAO, grpc.NewRewardServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
