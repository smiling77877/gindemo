package ioc

import (
	"gindemo/webbook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	// viper 设置默认值的方法
	// 1 使用viper的SetDefault的方法
	// 2 利用结构体, 在调用UnmarshalKey之前，设置好默认值
	var cfg Config = Config{
		DSN: viper.GetString("root:root@tcp(localhost:3316)/webook"),
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
