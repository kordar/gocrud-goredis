module github.com/kordar/gocrud-redis

go 1.18

replace github.com/kordar/gocrud => ../../github.com/gocrud

require (
	github.com/kordar/goutil v1.1.1
	github.com/redis/go-redis/v9 v9.6.1
	github.com/spf13/cast v1.7.0
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.3.0 // indirect
)
