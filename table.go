package gocrud_redis

type Table interface {
	TableName() string
	Find(v interface{}) Table
	FindOne(v interface{}) Table
	Count() int64
	Where(where Where) Table
	Match(value string) Table
	Or() Table
	Err() error
	Limit(offset int, limit int) Table
	Sort(key string, value string) Table
	Save(v interface{}) Table
}
