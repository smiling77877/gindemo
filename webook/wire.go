//go:build wireinject

package main

import (
	"gindemo/webook/interactive/events"
	repository2 "gindemo/webook/interactive/repository"
	cache2 "gindemo/webook/interactive/repository/cache"
	dao2 "gindemo/webook/interactive/repository/dao"
	service2 "gindemo/webook/interactive/service"
	"gindemo/webook/internal/events/article"
	"gindemo/webook/internal/repository"
	"gindemo/webook/internal/repository/cache"
	"gindemo/webook/internal/repository/dao"
	"gindemo/webook/internal/service"
	"gindemo/webook/internal/web"
	webook "gindemo/webook/internal/web/jwt"
	"gindemo/webook/ioc"
	"github.com/google/wire"
)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

var rankingSvcSet = wire.NewSet(
	cache.NewRankingRedisCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		ioc.InitSaramaClient,
		ioc.InitSyncProducer,
		ioc.InitRlockClient,
		// DAO部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,

		interactiveSvcSet,
		ioc.InitIntrClient,
		rankingSvcSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		article.NewSaramaSyncProducer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,

		// cache部分
		cache.NewCodeCache, cache.NewUserCache,
		cache.NewArticleRedisCache,
		// repository部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,
		repository.NewCachedArticleRepository,

		//
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// handler部分
		web.NewUserHandler,
		web.NewArticleHandler,
		webook.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
