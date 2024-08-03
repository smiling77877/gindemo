//go:build wireinject

package startup

import (
	repository2 "gindemo/webbook/interactive/repository"
	cache2 "gindemo/webbook/interactive/repository/cache"
	dao2 "gindemo/webbook/interactive/repository/dao"
	service2 "gindemo/webbook/interactive/service"
	"gindemo/webbook/internal/events/article"
	"gindemo/webbook/internal/job"
	"gindemo/webbook/internal/repository"
	"gindemo/webbook/internal/repository/cache"
	"gindemo/webbook/internal/repository/dao"
	"gindemo/webbook/internal/service"
	"gindemo/webbook/internal/service/sms"
	"gindemo/webbook/internal/service/sms/async"
	"gindemo/webbook/internal/web"
	ijwt "gindemo/webbook/internal/web/jwt"
	"gindemo/webbook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB, InitSaramaClient, InitSyncProducer, InitLogger)

var jobProviderSet = wire.NewSet(
	service.NewCronJobService,
	repository.NewPreemptJobRepository,
	dao.NewGORMJobDAO,
)

var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService)

var articleSvcProvider = wire.NewSet(
	repository.NewCachedArticleRepository,
	cache.NewArticleRedisCache,
	dao.NewArticleGORMDAO,
	service.NewArticleService)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		articleSvcProvider,
		interactiveSvcSet,
		// cache部分
		cache.NewCodeCache,
		// repository部分
		repository.NewCodeRepository,

		article.NewSaramaSyncProducer,

		// Service部分
		ioc.InitSMSService,
		service.NewCodeService,
		InitWechatService,
		// handler部分
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitAsyncSmsService(svc sms.Service) *async.Service {
	wire.Build(thirdPartySet, repository.NewAsyncSmsRepository,
		dao.NewGORMAsyncSmsDAO, async.NewService)
	return &async.Service{}
}

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		interactiveSvcSet,
		repository.NewCachedArticleRepository,
		cache.NewArticleRedisCache,
		service.NewArticleService,
		article.NewSaramaSyncProducer,
		web.NewArticleHandler)
	return &web.ArticleHandler{}
}

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdPartySet, interactiveSvcSet)
	return service2.NewInteractiveService(nil)
}

func InitJobScheduler() *job.Scheduler {
	wire.Build(jobProviderSet, thirdPartySet, job.NewScheduler)
	return &job.Scheduler{}
}
