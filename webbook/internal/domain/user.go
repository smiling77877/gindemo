package domain

import (
	"time"
)

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string

	Phone string

	Ctime time.Time
}

// TodayIsBirthday 判定今天是不是我的生日
func (u User) TodayIsBirthday() bool {
	now := time.Now()
	return now.Month() == u.Birthday.Month() && now.Day() == u.Birthday.Day()
}
