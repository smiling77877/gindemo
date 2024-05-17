package main

import (
	"gindemo/webbook/internal/repository"
	"gindemo/webbook/internal/repository/dao"
	"gindemo/webbook/internal/service"
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
	server := initWebServer()
	initUserHdl(db, server)
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
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
		Addr: "localhost:6379",
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
