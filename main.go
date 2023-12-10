package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	server := initWeebServer()

	initUserHdl(db, server)

	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uhdl := web.NewUserHandler(us)

	uhdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db

}

func initWeebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{

		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
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

	login := &middleware.LoginMiddlewareBuilder{}
	//存储数据，将userId存入Cookie中
	store := cookie.NewStore([]byte("secret"))

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
