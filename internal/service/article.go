package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository

	//V1写法
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	logger     logger.LoggerV1
}

func NewArticleServiceV1(readerRepo repository.ArticleReaderRepository, authorRepo repository.ArticleAuthorRepository, logger logger.LoggerV1) *articleService {
	return &articleService{
		readerRepo: readerRepo,
		authorRepo: authorRepo,
		logger:     logger,
	}
}
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}
func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	} else {
		return a.repo.Create(ctx, art)
	}
	return a.repo.Create(ctx, art)
}
func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	return a.repo.Sync(ctx, art)
}
func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	//先操作制作库，再操作线上库
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}

	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			a.logger.Error("保存到制作库成功但是到线上库失败", logger.Int64("aid", art.Id), logger.Error(err))
		} else {
			return id, nil
		}
	}
	a.logger.Error("保存到制作库成功但是到线上库失败，重试次数耗尽", logger.Int64("aid", art.Id), logger.Error(err))
	return id, errors.New("保存到线上库失败，重试次数耗尽")
}
