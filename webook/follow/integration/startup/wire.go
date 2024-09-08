//go:build wireinject

package startup

import (
	"gindemo/webook/follow/grpc"
	"gindemo/webook/follow/repository"
	"gindemo/webook/follow/repository/cache"
	"gindemo/webook/follow/repository/dao"
	"gindemo/webook/follow/service"
	"github.com/google/wire"
)

func InitServer() *grpc.FollowServiceServer {
	wire.Build(
		InitRedis, InitLog, InitTestDB, dao.NewGORMFollowRelationDAO,
		cache.NewRedisFollowCache, repository.NewFollowRelationRepository,
		service.NewFollowRelationService, grpc.NewFollowRelationServiceServer,
	)
	return new(grpc.FollowServiceServer)
}
