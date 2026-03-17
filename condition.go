package gocrud_redis

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/spf13/cast"
)

func asTable(db interface{}) (Table, bool) {
	if db == nil {
		return nil, false
	}
	t, ok := db.(Table)
	return t, ok
}

type whereFunc struct {
	key  string
	exec func(item map[string]interface{}) bool
}

func (w whereFunc) Key() string {
	return w.key
}

func (w whereFunc) Exec(item map[string]interface{}) bool {
	return w.exec(item)
}

func EQ(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			return equalValue(item[field], value)
		},
	})
}

func NEQ(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			return !equalValue(item[field], value)
		},
	})
}

func LT(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	return compareExecutor(db, field, value, func(c int) bool { return c < 0 })
}

func LE(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	return compareExecutor(db, field, value, func(c int) bool { return c <= 0 })
}

func GT(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	return compareExecutor(db, field, value, func(c int) bool { return c > 0 })
}

func GE(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	return compareExecutor(db, field, value, func(c int) bool { return c >= 0 })
}

func IN(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	values := toSlice(value)
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			v := item[field]
			for _, vv := range values {
				if equalValue(v, vv) {
					return true
				}
			}
			return false
		},
	})
}

func NOTIN(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	values := toSlice(value)
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			v := item[field]
			for _, vv := range values {
				if equalValue(v, vv) {
					return false
				}
			}
			return true
		},
	})
}

func LIKE(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	needle := cast.ToString(value)
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			hay := cast.ToString(item[field])
			return strings.Contains(hay, needle)
		},
	})
}

func NOTLIKE(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	needle := cast.ToString(value)
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			hay := cast.ToString(item[field])
			return !strings.Contains(hay, needle)
		},
	})
}

func LIKELEFT(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	prefix := cast.ToString(value)
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			hay := cast.ToString(item[field])
			return strings.HasPrefix(hay, prefix)
		},
	})
}

func LIKERIGHT(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	suffix := cast.ToString(value)
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			hay := cast.ToString(item[field])
			return strings.HasSuffix(hay, suffix)
		},
	})
}

func BETWEEN(db interface{}, field string, value interface{}, value2 ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	max := interface{}(nil)
	if len(value2) > 0 {
		max = value2[0]
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			v := item[field]
			return compare(v, value) >= 0 && compare(v, max) <= 0
		},
	})
}

func NOTBETWEEN(db interface{}, field string, value interface{}, value2 ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	max := interface{}(nil)
	if len(value2) > 0 {
		max = value2[0]
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			v := item[field]
			return !(compare(v, value) >= 0 && compare(v, max) <= 0)
		},
	})
}

func ISNULL(db interface{}, field string, _ interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			v, exists := item[field]
			if !exists || v == nil {
				return true
			}
			return strings.TrimSpace(cast.ToString(v)) == ""
		},
	})
}

func ISNOTNULL(db interface{}, field string, _ interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			v, exists := item[field]
			if !exists || v == nil {
				return false
			}
			return strings.TrimSpace(cast.ToString(v)) != ""
		},
	})
}

func ASC(db interface{}, field string, _ interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Sort(field, "ASC")
}

func DESC(db interface{}, field string, _ interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Sort(field, "DESC")
}

func PAGE(db interface{}, _ string, value interface{}, value2 ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	offset := cast.ToInt(value)
	limit := 0
	if len(value2) > 0 {
		limit = cast.ToInt(value2[0])
	}
	return t.Limit(offset, limit)
}

func CREATE(db interface{}, _ string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return errInvalidDB
	}
	t.Save(value)
	return t.Err()
}

func SAVE(db interface{}, _ string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return errInvalidDB
	}
	t.Save(value)
	return t.Err()
}

func SETVAL(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return errInvalidDB
	}
	t.Updates(map[string]interface{}{field: value})
	return t.Err()
}

func UPDATE(db interface{}, field string, value interface{}, _ ...interface{}) interface{} {
	return SETVAL(db, field, value)
}

func UPDATES(db interface{}, _ string, value interface{}, value2 ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return errInvalidDB
	}
	var data map[string]interface{}
	if len(value2) == 0 {
		data = toStringAnyMap(value)
	} else {
		data = toStringAnyMap(value2[0])
	}
	t.Updates(data)
	return t.Err()
}

func DELETE(db interface{}, _ string, _ interface{}, _ ...interface{}) interface{} {
	t, ok := asTable(db)
	if !ok {
		return errInvalidDB
	}
	t.Delete()
	return t.Err()
}

var errInvalidDB = errors.New("invalid db")

func compareExecutor(db interface{}, field string, value interface{}, okFn func(c int) bool) interface{} {
	t, ok := asTable(db)
	if !ok {
		return db
	}
	return t.Where(whereFunc{
		key: field,
		exec: func(item map[string]interface{}) bool {
			return okFn(compare(item[field], value))
		},
	})
}

func equalValue(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	af, aok := asFloat64(a)
	bf, bok := asFloat64(b)
	if aok && bok {
		return af == bf
	}
	return cast.ToString(a) == cast.ToString(b)
}

func compare(a interface{}, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	af, aok := asFloat64(a)
	bf, bok := asFloat64(b)
	if aok && bok {
		if af < bf {
			return -1
		}
		if af > bf {
			return 1
		}
		return 0
	}
	return strings.Compare(cast.ToString(a), cast.ToString(b))
}

func asFloat64(v interface{}) (float64, bool) {
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
		s := strings.TrimSpace(cast.ToString(v))
		if s == "" {
			return 0, false
		}
		f := cast.ToFloat64(v)
		return f, true
	}
}

func toSlice(v interface{}) []interface{} {
	switch t := v.(type) {
	case []interface{}:
		return t
	case []string:
		out := make([]interface{}, 0, len(t))
		for _, s := range t {
			out = append(out, s)
		}
		return out
	case []int:
		out := make([]interface{}, 0, len(t))
		for _, s := range t {
			out = append(out, s)
		}
		return out
	case []int64:
		out := make([]interface{}, 0, len(t))
		for _, s := range t {
			out = append(out, s)
		}
		return out
	default:
		s := strings.TrimSpace(cast.ToString(v))
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		out := make([]interface{}, 0, len(parts))
		for _, p := range parts {
			out = append(out, strings.TrimSpace(p))
		}
		return out
	}
}

func toStringAnyMap(v interface{}) map[string]interface{} {
	if v == nil {
		return map[string]interface{}{}
	}
	switch t := v.(type) {
	case map[string]interface{}:
		return t
	case map[string]string:
		out := map[string]interface{}{}
		for k, v := range t {
			out[k] = v
		}
		return out
	default:
		bs, err := json.Marshal(v)
		if err != nil {
			return map[string]interface{}{}
		}
		out := map[string]interface{}{}
		_ = json.Unmarshal(bs, &out)
		return out
	}
}
