package dao

import (
	"context"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./article_author.go -package=daomocks -destination=./mocks/article_author.mock.go ArticleAuthorDAO
type ArticleAuthorDAO interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type ArticleGORMAuthorDAO struct {
	db *gorm.DB
}

func (a *ArticleGORMAuthorDAO) Create(ctx context.Context, art Article) (int64, error) {
	panic("1")
}

func (a *ArticleGORMAuthorDAO) Update(ctx context.Context, art Article) error {
	panic("1")
}

func NewArticleGORMAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &ArticleGORMAuthorDAO{db: db}
}
