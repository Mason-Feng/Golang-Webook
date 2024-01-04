package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}
func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{cors.New(cors.Config{

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
	}),
		func(ctx *gin.Context) {
			println("这是我的Middleware")
		},
		ratelimit.NewBuilder(redisClient, time.Second, 1000).Build(),
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}
