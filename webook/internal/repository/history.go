package repository

import (
	"context"
	"gindemo/webook/internal/domain"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
