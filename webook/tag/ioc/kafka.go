package ioc

import (
	"gindemo/webook/tag/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.SyncProducer {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	producer, err := sarama.NewSyncProducer(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return producer
}

func InitProducer(client sarama.SyncProducer) events.Producer {
	res, err := events.NewSaramaSyncProducer(client)
	if err != nil {
		panic(err)
	}
	return res
}
