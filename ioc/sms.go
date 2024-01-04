package ioc

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
	"webook/internal/service/sms/tencent"
)

func InitSMSService() sms.SMSService {
	return localsms.NewService()
	//如果有需要腾讯的，就可以初始化
	//return initTencentSMSService()
}

func initTencentSMSService() sms.SMSService {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯SMS的secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯SMS的secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "xxxxx", "江苏大学")
}
