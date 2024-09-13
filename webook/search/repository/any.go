package repository

import (
	"context"
	"gindemo/webook/search/repository/dao"
)

type AnyRepository interface {
	Input(ctx context.Context, index, docID, data string) error
}

type anyRepository struct {
	dao dao.AnyDAO
}

func NewAnyRepository(dao dao.AnyDAO) AnyRepository {
	return &anyRepository{
		dao: dao,
	}
}

func (r *anyRepository) Input(ctx context.Context, index, docID, data string) error {
	return r.dao.Input(ctx, index, docID, data)
}
