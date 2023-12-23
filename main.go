package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"net/http"
	"webook/config"

	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()

	initUserHdl(db, server)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,启动成功了！")
	})

	server.Run(":8081")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uhdl := web.NewUserHandler(us)

	uhdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db

}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{

		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization", "User-Agent"},
		//允许前端访问后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
		//AllowOrigins:     []string{"*"},
		//AllowMethods:     []string{"GET", "POST", "DELETE", "HEAD", "OPTIONS", "PUT", "PATCH"},
		//AllowHeaders:     []string{"Origin"},
		//ExposeHeaders:    []string{"Content-Length"},
		//AllowCredentials: true,
		////AllowOriginFunc: func(origin string) bool {
		////	return origin == "https://github.com"
		////},
		//MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的Middleware")
	})
	//限流
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1).Build())

	useJWT(server)
	return server
}

func useJWT(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}
func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	//存储数据，将userId存入Cookie中
	store := cookie.NewStore([]byte("secret"))
	//基于内存的实现，第一个参数authentication key,最好是32或者64位
	//第二个参数是encryption key
	//store := memstore.NewStore([]byte("RrRqvf6sVUhBwm0hTl9Umu1vu1unNkp6"),
	//	[]byte("yZ3wbxxK28z67vLz0TiiY6br70mXFiHc"))
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("RrRqvf6sVUhBwm0hTl9Umu1vu1unNkp6"),
	//	[]byte("yZ3wbxxK28z67vLz0TiiY6br70mXFiHc"))
	//if err != nil {
	//	panic(err)
	//}

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
