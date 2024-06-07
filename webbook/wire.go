//go:build wireinject

package main

import (
	"gindemo/webbook/internal/repository"
	"gindemo/webbook/internal/repository/cache"
	"gindemo/webbook/internal/repository/dao"
	"gindemo/webbook/internal/service"
	"gindemo/webbook/internal/web"
	webook "gindemo/webbook/internal/web/jwt"
	"gindemo/webbook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// DAO部分
		dao.NewUserDAO,
		// cache部分
		cache.NewCodeCache, cache.NewUserCache,
		// repository部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,
		//
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		// handler部分
		web.NewUserHandler,
		webook.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
