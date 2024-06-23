package startup

import (
	"gindemo/webbook/internal/service/oauth2/wechat"
	"gindemo/webbook/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", nil)
}
