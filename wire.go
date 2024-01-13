//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"

	"webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//由下至上的顺序
		//第三方依赖
		ioc.InitRedis, ioc.InitDB,

		//DAO部分
		dao.NewUserDAO,

		//cache部分
		cache.NewRedisCodeCache, cache.NewUserCache,

		//repository部分
		repository.NewCacheUserRepository, repository.NewCodeRepository,

		//service部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,

		//handler部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
