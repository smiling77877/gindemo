package startup

import (
	dao2 "gindemo/webbook/interactive/repository/dao"
	"gindemo/webbook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	err = dao2.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
