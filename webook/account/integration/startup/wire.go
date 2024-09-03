//go:build wireinject

package startup

import (
	"gindemo/webook/account/grpc"
	"gindemo/webook/account/repository"
	"gindemo/webook/account/repository/dao"
	"gindemo/webook/account/service"
	"github.com/google/wire"
)

func InitAccountService() *grpc.AccountServiceServer {
	wire.Build(InitTestDB, dao.NewCreditGORMDAO, repository.NewAccountRepository,
		service.NewAccountService, grpc.NewAccountServiceServer)
	return new(grpc.AccountServiceServer)
}
