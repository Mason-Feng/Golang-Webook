package repository

import (
	"context"
	"gorm.io/gorm"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao       dao.ArticleDAO
	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO
	db        *gorm.DB
}

func NewCacheArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}
func NewCacheArticleRepositoryV2(readerDAO dao.ArticleReaderDAO, authorDAO dao.ArticleAuthorDAO) ArticleRepository {
	return &CacheArticleRepository{
		readerDAO: readerDAO,
		authorDAO: authorDAO,
	}
}

func NewCacheArticleRepositoryV3(readerDAO dao.ArticleReaderDAO, authorDAO dao.ArticleAuthorDAO) *CacheArticleRepository {
	return &CacheArticleRepository{
		readerDAO: readerDAO,
		authorDAO: authorDAO,
	}
}
func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(art))
}
func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))
}
func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error

	}
	//防止后面业务panic
	defer tx.Rollback()
	authorDAO := dao.NewArticleGORMAuthorDAO(tx)
	readerDAO := dao.NewArticleGORMReaderDAO(tx)
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = authorDAO.Update(ctx, artn)
	} else {
		id, err = authorDAO.Create(ctx, artn)
	}

	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle(artn))
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, err

}
func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = c.authorDAO.Update(ctx, artn)
	} else {
		id, err = c.authorDAO.Create(ctx, artn)
	}

	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = c.readerDAO.Upsert(ctx, artn)
	return id, err

}
