package gocrud_redis

import (
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

var (
	data = map[string]map[string]fieldCache{}
)

type fieldCache struct {
	Key       string
	FieldName string
	FieldType interface{}
	Index     int
}

func GetCache(table string, obj interface{}) map[string]fieldCache {
	m := data[table]
	if m != nil {
		return m
	}

	m = map[string]fieldCache{}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		// 遍历结构体字段
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i) // 获取第 i 个字段
			name := field.Name  // 字段名
			// 获取字段标签
			jsonTag := field.Tag.Get("json")
			m[jsonTag] = fieldCache{
				Key:       jsonTag,
				FieldName: name,
				FieldType: field.Type.Name(),
				Index:     i,
			}
		}
		data[table] = m
		return m
	default:
		return nil
	}

}

func GetCacheFieldId(table string, v interface{}, ids []string) string {
	cache := GetCache(table, v)
	names := make([]string, 0)
	if cache == nil {
		m := cast.ToStringMapString(v)
		for _, id := range ids {
			if c, ok := m[id]; ok {
				names = append(names, c)
			}
		}
	} else {
		value := reflect.ValueOf(v).Elem()
		for _, id := range ids {
			if c, ok := cache[id]; ok {
				field := value.Field(c.Index)
				names = append(names, field.String())
			}
		}
	}

	return strings.Join(names, "-")
}
