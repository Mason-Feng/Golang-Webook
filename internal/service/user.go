package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"webook/internal/domain"
	"webook/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}

}
func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.repo.Create(ctx, u)

}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateNonSensitiveInfo(ctx, u)

}

func (svc *UserService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userId)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		//有两种情况
		//err==nil，u是可用的
		//err!=nil，系统错误
		return u, err
	}
	//用户没找到，进行注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	//有两种可能，一种是err恰好是唯一索引冲突（phone)，
	//一种是err!=nil,系统错误
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	//要么err==nil，要么ErrDuplicateUser也代表用户存在
	//这里会出现主从延迟，理论上讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
