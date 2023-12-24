package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//带有sql.NullString标签的列表示这是一个可以为NULL的列
	Email sql.NullString `gorm:"unique"`
	Phone sql.NullString `gorm:"unique"`

	Password string
	Nickname string `gorm:"type=varchar(128)"`
	Brithday string
	AboutMe  string `gorm:"type=varchar(4096)"`
	Ctime    int64  // 创建时间,时区 UTC 0毫秒数
	Utime    int64  // 更新时间

}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error

	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			//用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}

	}
	return err
}

func (dao *UserDAO) UpdateNonSensitiveInfo(ctx context.Context, u User) error {
	user, err := dao.FindById(ctx, u.Id)
	if err != nil {
		return err
	}
	user.Nickname = u.Nickname
	user.Brithday = u.Brithday
	user.AboutMe = u.AboutMe
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err = dao.db.WithContext(ctx).Save(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			//用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}

	}

	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, userId int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("Id=?", userId).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("Phone=?", phone).First(&u).Error
	return u, err
}
