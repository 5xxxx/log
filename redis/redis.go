package redis

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

func NewRedisWriter(ctx context.Context, key string, cli *redis.Client) *redisWriter {
	return &redisWriter{
		cli: cli, listKey: key, ctx: ctx,
	}
}

// 为 logger 提供写入 redis 队列的 io 接口
type redisWriter struct {
	cli     *redis.Client
	listKey string
	ctx     context.Context
	sync.Mutex
}

func (w *redisWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	n, err := w.cli.RPush(w.ctx, w.listKey, p).Result()
	return int(n), err
}
