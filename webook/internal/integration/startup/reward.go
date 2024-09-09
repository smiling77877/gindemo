package startup

import (
	rewardv1 "gindemo/webook/api/proto/gen/reward/v1"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitRewardServiceClient() rewardv1.RewardServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.reward", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.NewClient(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	res := rewardv1.NewRewardServiceClient(cc)
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = Config{}
		err = viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			// 这边做不了什么
			panic(err)
		}
	})
	return res
}
