package ioc

import (
	"gindemo/webbook/internal/web"
	webook "gindemo/webbook/internal/web/jwt"
	"gindemo/webbook/internal/web/middleware"
	"gindemo/webbook/pkg/ginx/middleware/ratelimit"
	"gindemo/webbook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, hdl webook.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token", "x-refresh_token"},
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost")
			},
			MaxAge: 12 * time.Hour,
		}),
		//ratelimit.NewBuilder(redisClient, time.Second, 1).Build(),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		(middleware.NewLoginJWTMiddlewareBuilder(hdl)).CheckLogin(),
	}
}
