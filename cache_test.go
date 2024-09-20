package gocrud_redis

import (
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

func client() redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"192.168.0.190:6379"},
		Password: "940430Dex",
		DB:       2,
	})
}

func TestName(t *testing.T) {
	client := client()
	demo := Demo{"tom", 82}
	table := NewHashTable(client, "table", "id", "name")
	table.Save(&demo)
	//table.Save(map[string]interface{}{"id": "cc", "name": "ppp", "vv": 99999999, "age": "ooooooooooo"})
}

func TestIndex(t *testing.T) {
	rdb := client()
	newvalue := "2024-03-02 11:10:04"

	key := "MMM"
	index := NewIndex(rdb, key, 2)
	index.AddList(newvalue)

}
