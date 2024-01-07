package limiter

import "context"

type Limiter interface {
	//是否触发限流
	//返回true,就是触发限流
	Limit(ctx context.Context, key string) (bool, error)
}
