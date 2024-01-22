package main

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	//db := initDB()
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server := initWebServer()
	//codeSvc := initCodeSvc(redisClient)
	//
	//initUserHdl(db, redisClient, codeSvc, server)
	//server := gin.Default()
	initViperRemote()
	//initViperWatch()
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,启动成功了！")
	})

	server.Run(":8080")
}
func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
func initViperWatch() {
	cfile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.Set("db.dsn", "root:root@tcp(localhost:13316)/webook")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println(viper.GetString("test.key"))
	})

}
func initViper() {

	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	//当前工作目录的config子目录
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperV1() {

	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/dev.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperV2() {
	cfg := `test:
  key: value1
redis:
  addr: "localhost:6379"

db:
  dsn: "root:root@tcp(localhost:13316)/webook"`

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}
func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置中心发生变更")
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("watch", viper.GetString("test.key"))
			time.Sleep(time.Second * 3)
		}
	}()

}

//
//func initUserHdl(db *gorm.DB, redisClient redis.Cmdable, codeSvc service.CodeService, server *gin.Engine) {
//	ud := dao.NewUserDAO(db)
//
//	uc := cache.NewUserCache(redisClient)
//	ur := repository.NewCacheUserRepository(ud, uc)
//	us := service.NewUserService(ur)
//	uhdl := web.NewUserHandler(us, codeSvc)
//
//	uhdl.RegisterRoutes(server)
//}
//
//func initCodeSvc(redisClient redis.Cmdable) *service.codeService {
//	cc := cache.NewCodeCache(redisClient)
//	crepo := repository.NewCodeRepository(cc)
//	return service.NewCodeService(crepo, initMemorySms())
//}
//func initMemorySms() sms.SMSService {
//	return localsms.NewService()
//}
//
//func initWebServer() *gin.Engine {
//	server := gin.Default()
//
//	//server.Use(cors.New(cors.Config{
//	//
//	//	AllowCredentials: true,
//	//	AllowHeaders:     []string{"Content-Type", "Authorization", "User-Agent"},
//	//	//允许前端访问后端响应中带的头部
//	//	ExposeHeaders: []string{"x-jwt-token"},
//	//	AllowOriginFunc: func(origin string) bool {
//	//		return true
//	//	},
//	//	MaxAge: 12 * time.Hour,
//	//	//AllowOrigins:     []string{"*"},
//	//	//AllowMethods:     []string{"GET", "POST", "DELETE", "HEAD", "OPTIONS", "PUT", "PATCH"},
//	//	//AllowHeaders:     []string{"Origin"},
//	//	//ExposeHeaders:    []string{"Content-Length"},
//	//	//AllowCredentials: true,
//	//	////AllowOriginFunc: func(origin string) bool {
//	//	////	return origin == "https://github.com"
//	//	////},
//	//	//MaxAge: 12 * time.Hour,
//	//}), func(ctx *gin.Context) {
//	//	println("这是我的Middleware")
//	//})
//	//限流
//	//redisClient := redis.NewClient(&redis.Options{
//	//	Addr: config.Config.Redis.Addr,
//	//})
//	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1).Build())
//
//	useJWT(server)
//	return server
//}
//
//func useJWT(server *gin.Engine) {
//	login := middleware.LoginJWTMiddlewareBuilder{}
//	server.Use(login.CheckLogin())
//}
//func useSession(server *gin.Engine) {
//	login := &middleware.LoginMiddlewareBuilder{}
//	//存储数据，将userId存入Cookie中
//	store := cookie.NewStore([]byte("secret"))
//	//基于内存的实现，第一个参数authentication key,最好是32或者64位
//	//第二个参数是encryption key
//	//store := memstore.NewStore([]byte("RrRqvf6sVUhBwm0hTl9Umu1vu1unNkp6"),
//	//	[]byte("yZ3wbxxK28z67vLz0TiiY6br70mXFiHc"))
//	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
//	//	[]byte("RrRqvf6sVUhBwm0hTl9Umu1vu1unNkp6"),
//	//	[]byte("yZ3wbxxK28z67vLz0TiiY6br70mXFiHc"))
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
