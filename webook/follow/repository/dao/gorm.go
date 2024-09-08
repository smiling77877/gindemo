package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMFollowRelationDAO struct {
	db *gorm.DB
}

func NewGORMFollowRelationDAO(db *gorm.DB) FollowRelationDAO {
	return &GORMFollowRelationDAO{db: db}
}

func (g *GORMFollowRelationDAO) FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? AND status = ?", follower, FollowRelationsStatusActive).
		Offset(int(offset)).Limit(int(limit)).Find(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) FollowRelationDetail(ctx context.Context, follower, followee int64) (FollowRelation, error) {
	var res FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? AND followee = ? AND status = ?",
		follower, followee, FollowRelationsStatusActive).First(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CreateFollowRelation(ctx context.Context, c FollowRelation) error {
	// 我也要保持 insert or update 语义
	now := time.Now().UnixMilli()
	c.Ctime = now
	c.Utime = now
	c.Status = FollowRelationsStatusActive
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			// 这代表的是关注了-取消了-再关注了
			"status": FollowRelationsStatusActive,
			"utime":  now,
		}),
	}).Create(&c).Error
	// 在这里更新 FollowStatics 的计数 (也是 upsert)
}

func (g *GORMFollowRelationDAO) UpdateStatus(ctx context.Context, follower, followee int64, status uint8) error {
	// 当前 status 就是 inactive 的呢？
	// 不需要多此一举去检测我这个数据在不在，状态对不对
	return g.db.WithContext(ctx).Where("follower = ? AND followee = ?", follower, followee).
		Updates(map[string]any{
			"status": status,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).Select("count(follower)").
		// 如果要是没有额外索引，不用怀疑，全表扫描
		// 可以考虑在 followee 额外创建一个索引
		Where("followee = ? AND status = ?", uid, FollowRelationsStatusActive).Count(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).Select("count(followee)").
		// <follower, followee>
		Where("follower = ? AND status = ?", uid, FollowRelationsStatusActive).Count(&res).Error
	return res, err
}
