//go:build wireinject

package follow

import (
	grpc2 "gindemo/webook/follow/grpc"
	"gindemo/webook/follow/ioc"
	"gindemo/webook/follow/repository"
	"gindemo/webook/follow/repository/cache"
	"gindemo/webook/follow/repository/dao"
	"gindemo/webook/follow/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	cache.NewRedisFollowCache,
	dao.NewGORMFollowRelationDAO, repository.NewFollowRelationRepository,
	service.NewFollowRelationService, grpc2.NewFollowRelationServiceServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitDB, ioc.InitLogger, ioc.InitRedis)

func Init() *App {
	wire.Build(
		thirdProvider, serviceProviderSet, ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"))
	return new(App)
}
