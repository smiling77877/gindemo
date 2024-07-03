//go:build wireinject

package startup

import (
	"gindemo/webbook/internal/repository"
	"gindemo/webbook/internal/repository/cache"
	"gindemo/webbook/internal/repository/dao"
	"gindemo/webbook/internal/service"
	"gindemo/webbook/internal/web"
	ijwt "gindemo/webbook/internal/web/jwt"
	"gindemo/webbook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB, InitLogger)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		// DAO部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,
		// cache部分
		cache.NewCodeCache, cache.NewUserCache,
		// repository部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,
		repository.NewCachedArticleRepository,

		// Service部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
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

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		service.NewArticleService,
		web.NewArticleHandler,
		repository.NewCachedArticleRepository,
	)
	return &web.ArticleHandler{}
}
