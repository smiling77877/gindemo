package repository

import (
	"context"
	"gindemo/webook/internal/domain"
	"gindemo/webook/internal/repository/dao"

	"github.com/ecodeclub/ekit/sqlx"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

//go:generate mockgen -source=./async_sms_repository.go -package=repomocks -destination=mocks/async_sms_repository.mock.go AsyncSmsRepository
type AsyncSmsRepository interface {
	// Add 添加一个异步SMS记录
	// 你叫做 Create 或者 Insert 也可以
	Add(ctx context.Context, s domain.AsyncSMS) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSMS, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
}

type asyncSmsRepository struct {
	dao dao.AsyncSMSDAO
}

func NewAsyncSmsRepository(dao dao.AsyncSMSDAO) AsyncSmsRepository {
	return &asyncSmsRepository{
		dao: dao,
	}
}

func (a *asyncSmsRepository) Add(ctx context.Context, s domain.AsyncSMS) error {
	return a.dao.Insert(ctx, dao.AsyncSMS{
		Config: sqlx.JsonColumn[dao.SmsConfig]{
			Val: dao.SmsConfig{
				TplId:   s.TplId,
				Args:    s.Args,
				Numbers: s.Numbers,
			},
			Valid: true,
		},
		RetryMax: s.RetryMax,
	})
}

func (a *asyncSmsRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSMS, error) {
	as, err := a.dao.GetWaitingSMS(ctx)
	if err != nil {
		return domain.AsyncSMS{}, err
	}
	return domain.AsyncSMS{
		Id:       as.Id,
		TplId:    as.Config.Val.TplId,
		Numbers:  as.Config.Val.Numbers,
		Args:     as.Config.Val.Args,
		RetryMax: as.RetryMax,
	}, nil
}

func (a *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	if success {
		return a.dao.MarkSuccess(ctx, id)
	}
	return a.dao.MarkFailed(ctx, id)
}
