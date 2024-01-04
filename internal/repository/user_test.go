package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	cachemocks "webook/internal/repository/cache/mocks"
	"webook/internal/repository/dao"
	daomocks "webook/internal/repository/dao/mocks"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO)
		ctx  context.Context
		uid  int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Birthday: "100",
						AboutMe:  "自我介绍",
						Phone: sql.NullString{
							String: "15112345678",
							Valid:  true,
						},
						Ctime: 101,
						Utime: 102,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Birthday: "100",
					AboutMe:  "自我介绍",
					Phone:    "15112345678",
					Ctime:    time.UnixMilli(101),
				}).Return(nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: "100",
				AboutMe:  "自我介绍",
				Phone:    "15112345678",
				Ctime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "123456",
						Birthday: "100",
						AboutMe:  "自我介绍",
						Phone:    "15112345678",
						Ctime:    time.UnixMilli(101),
					}, nil)

				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: "100",
				AboutMe:  "自我介绍",
				Phone:    "15112345678",
				Ctime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{}, dao.ErrRecordNotFound)

				return c, d
			},
			uid:      123,
			ctx:      context.Background(),
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "回写缓存失败",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Birthday: "100",
						AboutMe:  "自我介绍",
						Phone: sql.NullString{
							String: "15112345678",
							Valid:  true,
						},
						Ctime: 101,
						Utime: 102,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Birthday: "100",
					AboutMe:  "自我介绍",
					Phone:    "15112345678",
					Ctime:    time.UnixMilli(101),
				}).Return(errors.New("redis错误"))
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: "100",
				AboutMe:  "自我介绍",
				Phone:    "15112345678",
				Ctime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, ud := tc.mock(ctrl)
			svc := NewCacheUserRepository(ud, uc)
			user, err := svc.FindById(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
