package ioc

import (
	"gindemo/webook/internal/service/sms"
	"gindemo/webook/internal/service/sms/localsms"
	"gindemo/webook/internal/service/sms/tencent"
	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	return localsms.NewService()
	//return initTencentSMSService()
}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯SMS的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯SMS的 secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "xxx", "xx科技")
}
