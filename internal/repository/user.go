package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

}
func (repo *UserRepository) UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error {
	return repo.dao.UpdateNonSensitiveInfo(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Brithday: u.Brithday,
		AboutMe:  u.AboutMe,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil

}

func (repo *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain1(u), err
}
func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *UserRepository) toDomain1(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Nickname: u.Nickname,
		Brithday: u.Brithday,
		AboutMe:  u.AboutMe,
	}
}
