//go:build wireinject

package main

import (
	"gindemo/webook/payment/grpc"
	"gindemo/webook/payment/ioc"
	"gindemo/webook/payment/repository"
	"gindemo/webook/payment/repository/dao"
	"gindemo/webook/payment/web"
	"gindemo/webook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitEtcdClient,
		ioc.InitKafka,
		ioc.InitProducer,
		ioc.InitWechatClient,
		dao.NewPaymentGORMDAO,
		ioc.InitDB,
		repository.NewPaymentRepository,
		grpc.NewWechatServiceServer,
		ioc.InitWechatNativeService,
		ioc.InitWechatConfig,
		ioc.InitWechatNotifyHandler,
		ioc.InitGRPCServer,
		web.NewWechatHandler,
		ioc.InitGinServer,
		ioc.InitLogger,
		wire.Struct(new(wego.App), "WebServer", "GRPCServer"))
	return new(wego.App)
}
