package cache

import "github.com/redis/go-redis/v9"

type CodeCache struct{
	cmd redis.Cmdable
}
