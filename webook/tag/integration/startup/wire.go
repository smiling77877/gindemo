//go:build wireinject

package startup

import (
	"gindemo/webook/tag/events"
	"gindemo/webook/tag/grpc"
	"gindemo/webook/tag/repository/cache"
	"gindemo/webook/tag/repository/dao"
	"gindemo/webook/tag/service"
	"github.com/google/wire"
)

func InitGRPCService(p events.Producer) *grpc.TagServiceServer {
	wire.Build(InitTestDB, InitRedis, InitLog,
		dao.NewGORMTagDAO, InitRepository, cache.NewRedisTagCache,
		service.NewTagService, grpc.NewTagServiceServer,
	)
	return new(grpc.TagServiceServer)
}
