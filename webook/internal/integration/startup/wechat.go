package startup

import (
	"gindemo/webook/internal/service/oauth2/wechat"
	"gindemo/webook/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", nil)
}
