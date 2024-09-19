package gocrud_redis

import (
	"context"
	"encoding/json"
	"github.com/kordar/goutil"
	"github.com/redis/go-redis/v9"
)

type HashTable struct {
	tableName         string
	match             string
	condition         []WhereWrapper
	rdb               redis.UniversalClient
	conditionPosition int
	err               error
	id                []string
}

func NewHashTable(rdb redis.UniversalClient, name string, id ...string) HashTable {
	return HashTable{
		tableName:         name,
		condition:         make([]WhereWrapper, 0),
		match:             "*",
		conditionPosition: -1,
		rdb:               rdb,
		id:                id,
	}
}

func (h *HashTable) Or() Table {
	h.conditionPosition++
	h.condition = append(h.condition, NewWhereWrapper())
	return h
}

func (h *HashTable) TableName() string {
	return h.tableName
}

func (h *HashTable) Match(value string) Table {
	h.match = value
	return h
}

func (h *HashTable) Find(v interface{}) Table {
	ctx := context.Background()
	iter := h.rdb.HScan(ctx, h.tableName, 0, h.match, 0).Iterator()
	result := make([]map[string]interface{}, 0)
	for iter.Next(ctx) {
		key := iter.Val()
		if bytes, err := h.rdb.HGet(ctx, h.tableName, key).Bytes(); err == nil {
			item := map[string]interface{}{}
			if err := json.Unmarshal(bytes, &item); err != nil {
				continue
			}
			if !Conditions(h.condition, item) {
				continue
			}
			result = append(result, item)
		}
	}

	switch v.(type) {
	case struct{}:
		break
	default:
		v = &result
	}

	h.err = iter.Err()
	return h
}

func (h *HashTable) FindOne(v interface{}) Table {
	//TODO implement me
	panic("implement me")
}

func (h *HashTable) Err() error {
	return h.err
}

func (h *HashTable) Count() int64 {
	//TODO implement me
	panic("implement me")
}

func (h *HashTable) Where(where Where) Table {
	h.condition[h.conditionPosition].AddWhere(where)
	return h
}

func (h *HashTable) Limit(offset int, limit int) Table {
	//TODO implement me
	panic("implement me")
}

func (h *HashTable) Sort(key string, value string) Table {
	//TODO implement me
	panic("implement me")
}

func (h *HashTable) primaryKey(value interface{}) string {
	ctx := context.Background()
	idValue := GetCacheFieldId(h.tableName, value, h.id)
	primaryKey := h.tableName + "-PRIMARY-KEY"
	if pid, err := h.rdb.HGet(ctx, primaryKey, idValue).Result(); err == nil {
		return pid
	} else {
		uuid := goutil.UUID()
		h.rdb.HSet(ctx, primaryKey, idValue, uuid)
		return uuid
	}
}

func (h *HashTable) Save(value interface{}) Table {
	ctx := context.Background()
	primaryKey := h.primaryKey(value)
	if marshal, err := json.Marshal(value); err == nil {
		h.rdb.HSet(ctx, h.tableName, primaryKey, string(marshal))
	} else {
		h.err = err
	}
	return h
}
