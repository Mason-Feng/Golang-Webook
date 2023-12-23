package repository

import (
	"context"
	"log"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
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
func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return du, nil
	}
	//只要err不为nil，就要查询数据库
	//err有两种
	//1.key不存在，说明Redis是正常的
	//2.访问Redis有问题。可能是网络有问题，也可能是redis本身就崩溃了
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	//往数据库里面写入数据，采用异步的方式，以能够提高查询性能
	go func() {
		err := repo.cache.Set(
			ctx,
			du,
		)

		if err != nil {
			log.Println(err)
			return
		}
	}()
	return du, nil
}
func (repo *UserRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	//只要err为nil，就返回
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		u, err := repo.dao.FindById(ctx, uid)

		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		err = repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}
		return du, nil
	default:
		return domain.User{}, err
	}
	//只要err不为nil，就要查询数据库
	//err有两种
	//1.key不存在，说明Redis是正常的
	//2.访问Redis有问题。可能是网络有问题，也可能是redis本身就崩溃了

	////往数据库里面写入数据，采用异步的方式，以能够提高查询性能
	//go func(){
	//	repo.cache.Set(ctx, du)
	//	if err!=nil{
	//		log.Println(err)
	//	}
	//}()

}
func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Brithday: u.Brithday,
	}
}

//func (repo *UserRepository) toDomain1(u dao.User) domain.User {
//	return domain.User{
//		Id:       u.Id,
//		Email:    u.Email,
//		Nickname: u.Nickname,
//		Birthday: u.Birthday,
//		AboutMe:  u.AboutMe,
//	}
//}
