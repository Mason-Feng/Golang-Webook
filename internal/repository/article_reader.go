package repository

import (
	"context"
	"webook/internal/domain"
)

type ArticleReaderRepository interface {
	//有则更新，无则插入
	Save(ctx context.Context, art domain.Article) error
}
