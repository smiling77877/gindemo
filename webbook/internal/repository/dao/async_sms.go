package dao

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
)

var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

type AsyncSMSDAO interface {
	Insert(ctx context.Context, s AsyncSms) error
	GetWaitingSMS(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

type AsyncSMS struct {
	id int64
	// 使用我在 ekit 里面支持的 JSON 字段
	Config sqlx.JsonColumn[SmsConfig]
	// 重试次数
	RetryCnt int
	// 重试的最大次数
	RetryMax int
	Status   uint8
	Ctime    int64
	Utime    int64 `gorm:"index"`
}

type SmsConfig struct {
	TplId   string
	Args    []string
	Numbers []string
}
