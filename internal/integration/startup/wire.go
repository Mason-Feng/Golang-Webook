//go:build wireinject

package startup

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

var thirdPartySet = wire.NewSet(
	InitRedis,
	InitDB,
	InitLogger,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//由下至上的顺序
		//第三方依赖
		thirdPartySet,

		//DAO部分
		dao.NewUserDAO,
		dao.NewArticleGORMDAO,

		//cache部分
		cache.NewRedisCodeCache, cache.NewUserCache,
		//repository部分
		repository.NewCacheUserRepository,
		repository.NewCodeRepository,
		repository.NewCacheArticleRepository,

		//service部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		InitWechatService,
		//handler部分
		web.NewUserHandler,
		web.NewArticleHandler,

		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,

		repository.NewCacheArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}

}
