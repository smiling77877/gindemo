//go:build wireinject

package startup

import (
	"gindemo/webbook/interactive/grpc"
	repository2 "gindemo/webbook/interactive/repository"
	cache2 "gindemo/webbook/interactive/repository/cache"
	dao2 "gindemo/webbook/interactive/repository/dao"
	service2 "gindemo/webbook/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB, InitLogger)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

func InitInteractiveService() *grpc.InteractiveServiceServer {
	wire.Build(thirdPartySet, interactiveSvcSet, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
