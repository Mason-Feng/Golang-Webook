package repository

import (
	"context"

	"webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
)

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: cache,
	}
}
func (c *CacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {

	err := c.cache.Set(ctx, biz, phone, code)
	return err
}
func (c *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
