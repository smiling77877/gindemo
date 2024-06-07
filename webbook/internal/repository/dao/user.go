package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound     = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateById(ctx context.Context, entity User) error
	FindById(ctx context.Context, uid int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openId).First(&u).Error
	return u, err
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) UpdateById(ctx context.Context, entity User) error {
	// 这种写法依赖于GORM的零值和主键更新特性
	// Update 非零值 WHERE id = ?
	// return dao.db.
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?",
		entity.Id).Updates(map[string]any{
		"utime":    time.Now().UnixMilli(),
		"nickname": entity.Nickname,
		"birthday": entity.Birthday,
		"about_me": entity.AboutMe,
	}).Error
}

func (dao *GORMUserDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&res).Error
	return res, err
}

type User struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 代表这是一个可以为NULL的列
	// Email *string
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	Phone sql.NullString

	// 1 如果查询要求同时使用 openid 和 unionid, 就要创建联合唯一索引
	// 2 如果查询只用 openid, 那么就在 openid 上创建唯一索引, 或者 <openid, unionid> 联合索引
	// 3 如果查询只用 unionid, 那么就在 unionid 上创建唯一索引, 或者 <union, openid> 联合索引
	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
