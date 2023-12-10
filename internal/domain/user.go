package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string
	Brithday string
	AboutMe  string
	Ctime    time.Time
}

//func (u User) ValidateEmail() bool{
//	//在这里用正则表达式校验
//	return u.Email
//}