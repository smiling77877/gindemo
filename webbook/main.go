package main

import (
	"gindemo/webbook/config"
	"gindemo/webbook/internal/repository"
	"gindemo/webbook/internal/repository/cache"
	"gindemo/webbook/internal/repository/dao"
	"gindemo/webbook/internal/service"
	"gindemo/webbook/internal/service/sms"
	"gindemo/webbook/internal/service/sms/localsms"
	"gindemo/webbook/internal/web"
	"gindemo/webbook/internal/web/middleware"
	"gindemo/webbook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server := initWebServer()
	codeSvc := initCodeSvc(redisClient)
	initUserHdl(db, redisClient, codeSvc, server)
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello, 启动成功了！")
	//})
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, redisClient redis.Cmdable, codeSvc *service.CodeService, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(redisClient)
	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us, codeSvc)
	hdl.RegisterRoutes(server)
}

func initCodeSvc(redisClient *redis.Client) *service.CodeService {
	cc := cache.NewCodeCache(redisClient)
	crepo := repository.NewCodeRepository(cc)
	return service.NewCodeService(crepo, initMemorySms())
}

func initMemorySms() sms.Service {
	return localsms.NewService()
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 只在初始化过程中panic
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1).Build())

	useJWT(server)
	// useSession(server)
	return server
}

func useJWT(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}

	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("kZ3cV0sR2aL5xX6yA8tR5wD7rF6lW1gU"),
	//	[]byte("pO5bK9qT4fJ0xK5nR0aA1wY2cK3yN1dG"))
	//if err != nil {
	//	panic(err)
	//}

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
