# gocrud-goredis

封装 `go-redis`，实现 [`gocrud`](https://github.com/kordar/gocrud) 的 Execute/Body 调用链，并提供基于 Redis Hash 的简单 Table 实现。

## 安装

```go
go get github.com/kordar/gocrud-redis@latest
```

## 初始化

```go
gocrud_redis.InitExec()
```

`InitExec()` 内部做了防重注册，多次调用是安全的。

## 快速使用（HashTable）

### 1) 创建 Table

```go
rdb := redis.NewUniversalClient(&redis.UniversalOptions{
	Addrs: []string{"127.0.0.1:6379"},
	DB:    0,
})

tbl := gocrud_redis.NewHashTable(rdb, "demo", "name")
table := &tbl
```

`NewHashTable(rdb, name, id...)` 的 `id...` 用来生成主键映射（复合主键会用 `-` 拼接）。

### 2) Create / Save

```go
gocrud_redis.InitExec()

type Demo struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

body := gocrud.NewFormBody("redis", context.Background())
body.Object = map[string]interface{}{
	"name": "tom",
	"age":  82,
}

form := gocrud_redis.NewRedisFormBody(body)
model := Demo{}
_, _ = form.Create(&model, table, nil)
```

### 3) Query / Paginate

```go
search := gocrud.NewSearchBody("redis", context.Background())
search.Page = 1
search.PageSize = 10
search.Conditions = []gocrud.Condition{
	{Key: "name", Type: "EQ", Value: "tom"},
}
search.Sorts = []gocrud.Sort{
	{Key: "age", Type: "DESC"},
}

sb := gocrud_redis.NewRedisSearchBody(search)
tx, _ := sb.RedisPaginate(table, nil)

list := make([]Demo, 0)
tx.Find(&list)
```

## 支持的 Condition.Type

- 比较：`EQ` `NEQ` `LT` `LE` `GT` `GE`
- 集合：`IN` `NOTIN`
- 模糊：`LIKE` `NOTLIKE` `LIKELEFT` `LIKERIGHT`
- 区间：`BETWEEN` `NOTBETWEEN`
- 空值：`ISNULL` `ISNOTNULL`
- 排序：`ASC` `DESC`
- 分页：内部换算 `offset=(page-1)*pageSize` 与 `limit=pageSize`

## 测试（可选）

需要 Redis 连接信息，通过环境变量注入：

```bash
set GOCRUD_REDIS_ADDRS=127.0.0.1:6379
set GOCRUD_REDIS_PASSWORD=
set GOCRUD_REDIS_DB=0
```
