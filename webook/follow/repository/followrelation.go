package repository

import (
	"context"
	"gindemo/webook/follow/domain"
	"gindemo/webook/follow/repository/cache"
	"gindemo/webook/follow/repository/dao"
	"gindemo/webook/pkg/logger"
)

type FollowRepository interface {
	// GetFollowee 获取某人的关注列表
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	// FollowInfo 查看关注人的详情
	FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error)
	// AddFollowRelation 创建关注关系
	AddFollowRelation(ctx context.Context, f domain.FollowRelation) error
	// InactiveFollowRelation 取消关注
	InactiveFollowRelation(ctx context.Context, follower, followee int64) error
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
}

func NewFollowRelationRepository(dao dao.FollowRelationDAO, cache cache.FollowCache, l logger.LoggerV1) FollowRepository {
	return &CachedRelationRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

type CachedRelationRepository struct {
	dao   dao.FollowRelationDAO
	cache cache.FollowCache
	l     logger.LoggerV1
}

func (d *CachedRelationRepository) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	// 你可以考虑在这里缓存关注者列表的第一页
	followerList, err := d.dao.FollowRelationList(ctx, follower, offset, limit)
	if err != nil {
		return nil, err
	}
	return d.genFollowRelationList(followerList), nil
}

func (d *CachedRelationRepository) genFollowRelationList(followerList []dao.FollowRelation) []domain.FollowRelation {
	res := make([]domain.FollowRelation, 0, len(followerList))
	for _, follower := range followerList {
		res = append(res, d.toDomain(follower))
	}
	return res
}

func (d *CachedRelationRepository) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	// follow:123:234 => 标签信息，分组信息
	c, err := d.dao.FollowRelationDetail(ctx, follower, followee)
	if err != nil {
		return domain.FollowRelation{}, err
	}
	return d.toDomain(c), err
}

func (d *CachedRelationRepository) AddFollowRelation(ctx context.Context, c domain.FollowRelation) error {
	err := d.dao.CreateFollowRelation(ctx, d.toEntity(c))
	if err != nil {
		return err
	}
	// 更新缓存里面的关注了多少人，以及有多少的粉丝的计数 +1
	return d.cache.Follow(ctx, c.Follower, c.Followee)
}

func (d *CachedRelationRepository) InactiveFollowRelation(ctx context.Context, follower, followee int64) error {
	err := d.dao.UpdateStatus(ctx, followee, follower, dao.FollowRelationsStatusInactive)
	if err != nil {
		return err
	}
	// -1
	return d.cache.CancelFollow(ctx, follower, followee)
}

// 获得个人的关注了多少人，以及粉丝的数量
func (d *CachedRelationRepository) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	res, err := d.cache.StaticsInfo(ctx, uid)
	if err == nil {
		return res, nil
	}
	// 没有我就去数据库里查询
	res.Followers, err = d.dao.CntFollower(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	res.Followees, err = d.dao.CntFollowee(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	err = d.cache.SetStaticsInfo(ctx, uid, res)
	if err != nil {
		// 这里记录一下日志就行
	}
	return res, nil
}

func (d *CachedRelationRepository) toDomain(r dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
		Followee: r.Followee,
		Follower: r.Follower,
	}
}

func (d *CachedRelationRepository) toEntity(c domain.FollowRelation) dao.FollowRelation {
	return dao.FollowRelation{
		Followee: c.Followee,
		Follower: c.Follower,
	}
}
