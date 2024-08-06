package dao

import (
	"context"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./article_reader.go -package=daomocks -destination=./mocks/article_reader.mock.go ArticleReaderDAO
type ArticleReaderDAO interface {
	// Upsert INSERT or UPDATE
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

type ArticleGORMReaderDAO struct {
	db *gorm.DB
}

func (a *ArticleGORMReaderDAO) Upsert(ctx context.Context, art Article) error {
	panic("implement me")
}

func (a *ArticleGORMReaderDAO) UpsertV2(ctx context.Context, art PublishedArticle) error {
	panic("implement me")
}

func NewArticleGORMReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &ArticleGORMReaderDAO{db: db}
}
