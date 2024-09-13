package events

import (
	"gindemo/webook/pkg/logger"
	"gindemo/webook/search/service"
	"github.com/IBM/sarama"
)

type SyncDataEvent struct {
	IndexName string
	DocID     string
	// 这里应该是 BizTags
	Data string
}

type SyncDataEventConsumer struct {
	svc    service.SyncService
	client sarama.Client
	l      logger.LoggerV1
}
