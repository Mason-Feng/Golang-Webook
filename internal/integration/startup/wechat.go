package startup

import (
	"webook/internal/service/oauth2/wechat"
	"webook/pkg/logger"
)

func InitWechatService(logger logger.LoggerV1) wechat.OAuth2Service {
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("找不到环境变量WECHAT_APP_ID")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("找不到环境变量WECHAT_APP_SECRET")
	//}

	appID := "wxbdc5610cc59c1631"
	appSecret := "12342"
	return wechat.NewOAuth2WechatService(appID, appSecret, logger)

}
