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
		ioc.InitLogger,

		//DAO部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,

		//cache部分
		cache.NewRedisCodeCache,
		cache.NewUserCache,

		//repository部分
		repository.NewCacheUserRepository,
		repository.NewCodeRepository,
		repository.NewCacheArticleRepository,

		//service部分
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		//handler部分
		web.NewUserHandler,
		web.NewArticleHandler,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
