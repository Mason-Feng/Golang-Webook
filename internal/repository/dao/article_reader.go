package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	//Upsert:INSERT or UPDATE
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

func (a *ArticleGORMReaderDAO) UpsertV2(ctx context.Context, art PublishedArticle) error {
	//TODO implement me
	panic("implement me")
}

type ArticleGORMReaderDAO struct {
	db *gorm.DB
}

func NewArticleGORMReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &ArticleGORMReaderDAO{db: db}
}
func (a *ArticleGORMReaderDAO) Upsert(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}
