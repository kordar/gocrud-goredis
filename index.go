package gocrud_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
)

type Index struct {
	rdb       redis.UniversalClient
	key       string
	rangeSize int64
}

func NewIndex(rdb redis.UniversalClient, key string, rangeSize int64) *Index {
	return &Index{rdb: rdb, key: key, rangeSize: rangeSize}
}

func (index Index) AddList(value string) {
	ctx := context.Background()
	size := index.rdb.LLen(ctx, index.key).Val()
	position, s := index.rangeInsertListPosition(size, value)
	if s == "before" {
		index.rdb.LInsertBefore(ctx, index.key, position, value)
	}
	if s == "after" {
		index.rdb.LInsertAfter(ctx, index.key, position, value)
	}
	if s == "none" {
		index.rdb.LPush(ctx, index.key, value)
	}
}

func (index Index) rangeInsertListPosition(size int64, value string) (string, string) {
	position := ""
	beforeOrAfter := "none"
	if size == 0 {
		return position, beforeOrAfter
	}
	ctx := context.Background()
	var i int64 = 0
	for ; i < size; i = i + index.rangeSize {
		end := i + index.rangeSize - 1
		if end >= size {
			end = size - 1
		}
		list, _ := index.rdb.LRange(ctx, index.key, i, end).Result()
		for _, item := range list {
			position = item
			if strings.Compare(value, item) > 0 {
				beforeOrAfter = "after"
			} else {
				return position, "before"
			}
		}
	}
	return position, beforeOrAfter
}
