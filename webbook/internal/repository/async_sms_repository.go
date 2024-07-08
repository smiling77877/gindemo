package repository

import (
	"context"
	"gindemo/webbook/internal/domain"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

type AsyncSmsRepository interface {
	// Add 添加一个异步SMS记录
	// 你叫做 Create 或者 Insert 也可以
	Add(ctx context.Context, s domain.AsyncSMS) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSMS, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
}
