//go:build wireinject

package startup

import (
	service2 "gindemo/webook/interactive/service"
	"gindemo/webook/internal/events/article"
	"gindemo/webook/internal/job"
	"gindemo/webook/internal/repository"
	"gindemo/webook/internal/repository/cache"
	"gindemo/webook/internal/repository/dao"
	"gindemo/webook/internal/service"
	"gindemo/webook/internal/service/sms"
	"gindemo/webook/internal/service/sms/async"
	"gindemo/webook/internal/web"
	ijwt "gindemo/webook/internal/web/jwt"
	"gindemo/webook/ioc"
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
