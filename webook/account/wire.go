//go:build wireinject

package account

import (
	"gindemo/webook/account/grpc"
	"gindemo/webook/account/ioc"
	"gindemo/webook/account/repository"
	"gindemo/webook/account/repository/dao"
	"gindemo/webook/account/service"
	"gindemo/webook/pkg/wego"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(ioc.InitDB, ioc.InitLogger, ioc.InitEtcdClient,
		ioc.InitGRPCxServer, dao.NewCreditGORMDAO,
		repository.NewAccountRepository, service.NewAccountService,
		grpc.NewAccountServiceServer, wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
