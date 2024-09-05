package service

import (
	"context"
	"gindemo/webook/comment/domain"
	"gindemo/webook/comment/repository"
)

type CommentService interface {
	// GetCommentList Comment的id为0 获取一级评论
	// 按照 ID 倒序排序
	GetCommentList(ctx context.Context, biz string, bizID, minID, limit int64) ([]domain.Comment, error)
	// DeleteComment 删除评论，删除本评论和其子评论
	DeleteComment(ctx context.Context, id int64) error
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	GetMoreReplies(ctx context.Context, rid, maxID, limit int64) ([]domain.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentSvc(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (c *commentService) GetCommentList(ctx context.Context, biz string, bizID, minID, limit int64) ([]domain.Comment, error) {
	list, err := c.repo.FindByBiz(ctx, biz, bizID, minID, limit)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *commentService) DeleteComment(ctx context.Context, id int64) error {
	return c.repo.DeleteComment(ctx, domain.Comment{
		Id: id,
	})
}

func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.repo.CreateComment(ctx, comment)
}

func (c *commentService) GetMoreReplies(ctx context.Context, rid, maxID, limit int64) ([]domain.Comment, error) {
	return c.repo.GetMoreReplies(ctx, rid, maxID, limit)
}
