package gocrud_redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/kordar/gocrud"
	"github.com/redis/go-redis/v9"
)

type Demo struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func TestGetCache(t *testing.T) {
	cache := GetCache("demo", Demo{"tom", 82})
	fmt.Printf("----------%v", cache)
}

func client(t *testing.T) redis.UniversalClient {
	t.Helper()
	addrs := strings.TrimSpace(os.Getenv("GOCRUD_REDIS_ADDRS"))
	if addrs == "" {
		t.Skip("GOCRUD_REDIS_ADDRS is empty")
	}
	db, _ := strconv.Atoi(strings.TrimSpace(os.Getenv("GOCRUD_REDIS_DB")))
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split(addrs, ","),
		Password: os.Getenv("GOCRUD_REDIS_PASSWORD"),
		DB:       db,
	})
}

func TestName(t *testing.T) {
	rdb := client(t)
	demo := Demo{"tom", 82}
	tbl := NewHashTable(rdb, "table", "id", "name")
	table := &tbl
	table.Save(&demo)
	//table.Save(map[string]interface{}{"id": "cc", "name": "ppp", "vv": 99999999, "age": "ooooooooooo"})
}

func TestIndex(t *testing.T) {
	rdb := client(t)
	newvalue := "2024-03-02 11:10:04"

	key := "MMM"
	index := NewIndex(rdb, key, 2)
	index.AddList(newvalue)

}

func TestGocrudIntegration(t *testing.T) {
	rdb := client(t)
	InitExec()

	tbl := NewHashTable(rdb, "demo", "name")
	table := &tbl

	formBody := gocrud.NewFormBody("redis", context.Background())
	formBody.Object = map[string]interface{}{
		"name": "tom",
		"age":  82,
	}
	form := NewRedisFormBody(formBody)
	model := Demo{}
	if _, err := form.Create(&model, table, nil); err != nil {
		t.Fatal(err)
	}

	searchBody := gocrud.NewSearchBody("redis", context.Background())
	searchBody.Conditions = []gocrud.Condition{
		{Key: "name", Type: "EQ", Value: "tom"},
	}
	search := NewRedisSearchBody(searchBody)
	tx := search.RedisQuery(table, nil)
	list := make([]Demo, 0)
	tx.Find(&list)
	if tx.Err() != nil {
		t.Fatal(tx.Err())
	}
}
