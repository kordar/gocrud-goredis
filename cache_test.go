package gocrud_redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

type Demo struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func TestGetCache(t *testing.T) {
	cache := GetCache("demo", Demo{"tom", 82})
	fmt.Printf("----------%v", cache)
}

func TestName(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"192.168.30.16:30202"},
		Password: "940430Dex",
		DB:       2,
	})
	ctx := context.Background()
	client.ZAdd(ctx, "AAAAAAAA", redis.Z{})
	//demo := Demo{"tom", 82}
	//table := NewHashTable(client, "table", "id", "name")
	//table.Save(&demo)
	//table.Save(map[string]interface{}{"id": "cc", "name": "ppp", "vv": 99999999, "age": "ooooooooooo"})
}
