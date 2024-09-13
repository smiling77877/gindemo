package events

import (
	"context"
	"gindemo/webook/pkg/logger"
	"gindemo/webook/pkg/saramax"
	"gindemo/webook/search/domain"
	"gindemo/webook/search/service"
	"github.com/IBM/sarama"
	"time"
)

const topicSyncUser = "sync_user_event"

type UserConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.LoggerV1
}

type UserEvent struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

func NewUserConsumer(client sarama.Client, l logger.LoggerV1, svc service.SyncService) *UserConsumer {
	return &UserConsumer{
		client:  client,
		l:       l,
		syncSvc: svc,
	}
}

func (u *UserConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_user", u.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{topicSyncUser}, saramax.NewHandler[UserEvent](u.l, u.Consume))
		if err != nil {
			u.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (u *UserConsumer) Consume(sg *sarama.ConsumerMessage, evt UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return u.syncSvc.InputUser(ctx, u.toDomain(evt))
}

func (u *UserConsumer) toDomain(user UserEvent) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Nickname: user.Nickname,
		Phone:    user.Phone,
	}
}
