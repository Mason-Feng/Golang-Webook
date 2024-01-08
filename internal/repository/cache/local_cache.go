package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"time"
)

type CodeLocalCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type LocalCache struct {
	lc *cache.Cache
}

func NewLocalCache(expiration time.Duration, cleanupinterval time.Duration) *LocalCache {

	return &LocalCache{
		lc: cache.New(expiration, cleanupinterval),
	}
}

type CodeItem struct {
	code string
	cnt  int64
}

func (lc *LocalCache) Set(ctx context.Context, biz, phone, code string) error {
	_, expiration, found := lc.lc.GetWithExpiration(lc.key(biz, phone))
	if !found {
		lc.lc.Set(lc.key(biz, phone), &CodeItem{
			code: code,
			cnt:  3,
		}, cache.DefaultExpiration)
		return nil
	}

	if expiration.Sub(time.Now()) > time.Minute*9 {
		return ErrCodeSendTooMany
	}
	lc.lc.Set(lc.key(biz, phone), &CodeItem{
		code: code,
		cnt:  3,
	}, cache.DefaultExpiration)
	return nil

}

func (lc *LocalCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, found := lc.lc.Get(lc.key(biz, phone))
	if !found {
		return false, ErrKeyNotExist
	}
	res_code, res_cnt := res.(*CodeItem).code, res.(*CodeItem).cnt //进行类型断言，并取出其中的值
	if res_cnt <= 0 {
		return false, ErrCodeVerifyTooMany
	}
	if res_code != code {
		res.(*CodeItem).cnt--
		return false, ErrKeyNotExist
	}
	return true, nil

}
func (lc *LocalCache) key(biz, phone string) string {
	return biz + phone
}
