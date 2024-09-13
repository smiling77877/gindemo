//go:build wireinject

package startup

import (
	"gindemo/webook/search/grpc"
	"gindemo/webook/search/ioc"
	"gindemo/webook/search/repository"
	"gindemo/webook/search/repository/dao"
	"gindemo/webook/search/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewUserElasticDAO, dao.NewArticleElasticDAO, dao.NewTagESDAO, dao.NewAnyESDAO,
	repository.NewUserRepository, repository.NewAnyRepository, repository.NewArticleRepository,
	service.NewSyncService, service.NewSearchService,
)

var thirdProvider = wire.NewSet(InitESClient, ioc.InitLogger)

func InitSearchServer() *grpc.SearchServiceServer {
	wire.Build(thirdProvider, serviceProviderSet, grpc.NewSearchService)
	return new(grpc.SearchServiceServer)
}

func InitSyncServer() *grpc.SyncServiceServer {
	wire.Build(thirdProvider, serviceProviderSet, grpc.NewSyncServiceServer)
	return new(grpc.SyncServiceServer)
}
