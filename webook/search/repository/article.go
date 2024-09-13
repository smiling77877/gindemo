package repository

import (
	"context"
	"gindemo/webook/search/domain"
	"gindemo/webook/search/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type articleRepository struct {
	dao  dao.ArticleDAO
	tags dao.TagDAO
}

func NewArticleRepository(dao dao.ArticleDAO, tags dao.TagDAO) ArticleRepository {
	return &articleRepository{
		dao:  dao,
		tags: tags,
	}
}

func (a *articleRepository) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error) {
	artIds, err := a.tags.Search(ctx, uid, "article", keywords)
	if err != nil {
		return nil, err
	}
	arts, err := a.dao.Search(ctx, artIds, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(arts, func(idx int, src dao.Article) domain.Article {
		return domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Status:  src.Status,
			Content: src.Content,
			Tags:    src.Tags,
		}
	}), nil
}

func (a *articleRepository) InputArticle(ctx context.Context, msg domain.Article) error {
	return a.dao.InputArticle(ctx, dao.Article{
		Id:      msg.Id,
		Title:   msg.Title,
		Status:  msg.Status,
		Content: msg.Content,
	})
}
