package repository

import (
	"context"

	"webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRespository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRespository {
	return &CodeRespository{
		cache: cache,
	}
}
func (c *CodeRespository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}
func (c *CodeRespository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
