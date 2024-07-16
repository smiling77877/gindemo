package ioc

import (
	"gindemo/webbook/internal/web"
	webook "gindemo/webbook/internal/web/jwt"
	"gindemo/webbook/internal/web/middleware"
	"gindemo/webbook/pkg/ginx/middleware/prometheus"
	"gindemo/webbook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	artHdl *web.ArticleHandler,
	wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, hdl webook.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		Namespace: "geektime_daming",
		Subsystem: "webook",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
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
		//ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		//middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
		//	l.Debug("", logger.Field{Key: "req", Val: al})
		//}).AllowReqBody().AllowRespBody().Build(),
		pb.BuildRepsponseTime(),
		pb.BuildActiveRequest(),
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
	}
}
