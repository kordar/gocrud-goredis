package gocrud_redis

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"

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
	offset            int
	limit             int
	sortKey           string
	sortOrder         string
}

func NewHashTable(rdb redis.UniversalClient, name string, id ...string) HashTable {
	return HashTable{
		tableName:         name,
		condition:         []WhereWrapper{NewWhereWrapper()},
		match:             "*",
		conditionPosition: 0,
		rdb:               rdb,
		id:                id,
		offset:            0,
		limit:             0,
		sortKey:           "",
		sortOrder:         "",
	}
}

func (h *HashTable) Or() Table {
	h.condition = append(h.condition, NewWhereWrapper())
	h.conditionPosition = len(h.condition) - 1
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
		primaryKey := iter.Val()
		if !iter.Next(ctx) {
			break
		}
		value := iter.Val()
		item := map[string]interface{}{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}
		if !Conditions(h.condition, item) {
			continue
		}
		item["_pk"] = primaryKey
		result = append(result, item)
	}

	if h.sortKey != "" {
		sort.SliceStable(result, func(i, j int) bool {
			a := result[i][h.sortKey]
			b := result[j][h.sortKey]
			if strings.EqualFold(h.sortOrder, "DESC") {
				return compareValue(a, b) > 0
			}
			return compareValue(a, b) < 0
		})
	}

	if h.offset < 0 {
		h.offset = 0
	}
	if h.limit > 0 {
		start := h.offset
		if start > len(result) {
			start = len(result)
		}
		end := start + h.limit
		if end > len(result) {
			end = len(result)
		}
		result = result[start:end]
	} else if h.offset > 0 {
		start := h.offset
		if start > len(result) {
			start = len(result)
		}
		result = result[start:]
	}

	if v != nil && reflect.ValueOf(v).Kind() == reflect.Ptr {
		bs, err := json.Marshal(result)
		if err != nil {
			h.err = err
		} else {
			if err := json.Unmarshal(bs, v); err != nil {
				h.err = err
			}
		}
	}

	h.err = iter.Err()
	return h
}

func (h *HashTable) FindOne(v interface{}) Table {
	originOffset, originLimit := h.offset, h.limit
	h.offset, h.limit = 0, 1
	list := make([]map[string]interface{}, 0)
	h.Find(&list)
	h.offset, h.limit = originOffset, originLimit
	if h.err != nil || v == nil {
		return h
	}
	if len(list) == 0 {
		return h
	}
	bs, err := json.Marshal(list[0])
	if err != nil {
		h.err = err
		return h
	}
	if err := json.Unmarshal(bs, v); err != nil {
		h.err = err
		return h
	}
	return h
}

func (h *HashTable) Err() error {
	return h.err
}

func (h *HashTable) Count() int64 {
	ctx := context.Background()
	iter := h.rdb.HScan(ctx, h.tableName, 0, h.match, 0).Iterator()
	var count int64 = 0
	for iter.Next(ctx) {
		_ = iter.Val()
		if !iter.Next(ctx) {
			break
		}
		value := iter.Val()
		item := map[string]interface{}{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}
		if !Conditions(h.condition, item) {
			continue
		}
		count++
	}
	h.err = iter.Err()
	return count
}

func (h *HashTable) Where(where Where) Table {
	if len(h.condition) == 0 {
		h.condition = []WhereWrapper{NewWhereWrapper()}
		h.conditionPosition = 0
	}
	if h.conditionPosition < 0 || h.conditionPosition >= len(h.condition) {
		h.conditionPosition = 0
	}
	h.condition[h.conditionPosition].AddWhere(where)
	return h
}

func (h *HashTable) Limit(offset int, limit int) Table {
	h.offset = offset
	h.limit = limit
	return h
}

func (h *HashTable) Sort(key string, value string) Table {
	h.sortKey = key
	h.sortOrder = value
	return h
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
		h.err = h.rdb.HSet(ctx, h.tableName, primaryKey, string(marshal)).Err()
	} else {
		h.err = err
	}
	return h
}

func (h *HashTable) Updates(data map[string]interface{}) Table {
	ctx := context.Background()
	iter := h.rdb.HScan(ctx, h.tableName, 0, h.match, 0).Iterator()
	for iter.Next(ctx) {
		primaryKey := iter.Val()
		if !iter.Next(ctx) {
			break
		}
		value := iter.Val()
		item := map[string]interface{}{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}
		if !Conditions(h.condition, item) {
			continue
		}
		for k, v := range data {
			item[k] = v
		}
		bs, err := json.Marshal(item)
		if err != nil {
			h.err = err
			break
		}
		if err := h.rdb.HSet(ctx, h.tableName, primaryKey, string(bs)).Err(); err != nil {
			h.err = err
			break
		}
	}
	if h.err == nil {
		h.err = iter.Err()
	}
	return h
}

func (h *HashTable) Delete() Table {
	ctx := context.Background()
	iter := h.rdb.HScan(ctx, h.tableName, 0, h.match, 0).Iterator()
	keys := make([]string, 0)
	for iter.Next(ctx) {
		primaryKey := iter.Val()
		if !iter.Next(ctx) {
			break
		}
		value := iter.Val()
		item := map[string]interface{}{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}
		if !Conditions(h.condition, item) {
			continue
		}
		keys = append(keys, primaryKey)
	}
	if err := iter.Err(); err != nil {
		h.err = err
		return h
	}
	if len(keys) == 0 {
		return h
	}
	h.err = h.rdb.HDel(ctx, h.tableName, keys...).Err()
	return h
}

func compareValue(a interface{}, b interface{}) int {
	af := strings.TrimSpace(castToString(a))
	bf := strings.TrimSpace(castToString(b))
	if af == "" && bf == "" {
		return 0
	}
	if af == "" {
		return -1
	}
	if bf == "" {
		return 1
	}
	ai, aok := castToFloat(a)
	bi, bok := castToFloat(b)
	if aok && bok {
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
		return 0
	}
	return strings.Compare(af, bf)
}

func castToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		bs, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		s := string(bs)
		return strings.Trim(s, "\"")
	}
}

func castToFloat(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	case uint:
		return float64(t), true
	case uint64:
		return float64(t), true
	case uint32:
		return float64(t), true
	default:
		s := castToString(v)
		if s == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
}
