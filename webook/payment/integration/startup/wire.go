//go:build wireinject

package startup

import (
	"gindemo/webook/payment/ioc"
	"gindemo/webook/payment/repository"
	"gindemo/webook/payment/repository/dao"
	"gindemo/webook/payment/service/wechat"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitLogger, InitTestDB)

var wechatNativeSvcSet = wire.NewSet(
	ioc.InitWechatClient,
	dao.NewPaymentGORMDAO,
	repository.NewPaymentRepository,
	ioc.InitWechatNativeService,
	ioc.InitWechatConfig,
)

func InitWechatNativeService() *wechat.NativePaymentService {
	wire.Build(wechatNativeSvcSet, thirdPartySet)
	return new(wechat.NativePaymentService)
}
