package domain

import "time"

type User struct {
	Id         int64
	Email      string
	Phone      string
	Password   string
	Nickname   string
	Birthday   string
	AboutMe    string
	Ctime      time.Time
	WechatInfo WechatInfo
}

//func (u User) ValidateEmail() bool{
//	//在这里用正则表达式校验
//	return u.Email
//}
