package startup

import "gindemo/webook/pkg/logger"

func InitLogger() logger.LoggerV1 {
	return logger.NewNopLogger()
}
