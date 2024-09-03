//go:build wireinject

package startup

import (
	accountv1 "gindemo/webook/api/proto/gen/account/v1"
	pmtv1 "gindemo/webook/api/proto/gen/payment/v1"
	"gindemo/webook/reward/ioc"
	"gindemo/webook/reward/repository"
	"gindemo/webook/reward/repository/cache"
	"gindemo/webook/reward/repository/dao"
	"gindemo/webook/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(InitTestDB, ioc.InitLogger, InitRedis)

func InitWechatNativeSvc(pclient pmtv1.WechatPaymentServiceClient, aclient accountv1.AccountServiceClient) *service.WechatNativeRewardService {
	wire.Build(service.NewWechatNativeRewardService,
		thirdPartySet, cache.NewRewardRedisCache, repository.NewRewardRepository, dao.NewRewardGORMDAO)
	return new(service.WechatNativeRewardService)
}
