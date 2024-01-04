package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository"
	repomocks "webook/internal/repository/mocks"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#fqw")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#fqw"))
	assert.NoError(t, err)

}

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "fqw@qq.com").Return(domain.User{
					Email:    "fqw@qq.com",
					Password: "$2a$10$nEf0f5q.UKLOhhHZi0ZyouPD5bW8iUQ7BvQjDrtw9mpTahWRlXTJC",
					Phone:    "18652878928",
				}, nil)
				return repo
			},
			email:    "fqw@qq.com",
			password: "123456#fqw",
			wantUser: domain.User{
				Email:    "fqw@qq.com",
				Password: "$2a$10$nEf0f5q.UKLOhhHZi0ZyouPD5bW8iUQ7BvQjDrtw9mpTahWRlXTJC",
				Phone:    "18652878928",
			},
		},
		{name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "fqw@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "fqw@qq.com",
			password: "123456#fqw",
			wantErr:  ErrInvalidUserOrPassword,
		},
		{name: "系统出错",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "fqw@qq.com").Return(domain.User{}, errors.New("db错误"))
				return repo
			},
			email:    "fqw@qq.com",
			password: "123456#fqw",
			wantErr:  errors.New("db错误"),
		},
		{name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "fqw@qq.com").Return(domain.User{
					Email:    "fqw@qq.com",
					Password: "$2a$10$nEf0f5q.UKLOhhHZi0ZyouPD5bW8iUQ7BvQjDrtw9mpTahWRlXTJC",
					Phone:    "18652878928",
				}, nil)
				return repo
			},
			email:    "fqw@qq.com",
			password: "123456#fqwQW",
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
