package gocrud_redis

type Where interface {
	Key() string
	Exec(item map[string]interface{}) bool
}

type WhereWrapper struct {
	items []Where
}

func NewWhereWrapper() WhereWrapper {
	return WhereWrapper{items: make([]Where, 0)}
}

func (w *WhereWrapper) AddWhere(where Where) {
	w.items = append(w.items, where)
}

func (w *WhereWrapper) Exec(value map[string]interface{}) bool {
	if len(w.items) == 0 {
		return false
	}
	for _, item := range w.items {
		if !item.Exec(value) {
			return false
		}
	}
	return true
}

func Conditions(wrappers []WhereWrapper, items map[string]interface{}) bool {
	for _, wrapper := range wrappers {
		if wrapper.Exec(items) {
			return true
		}
	}
	return false
}

//
//func NewWhere(key string, value ...interface{}) Where {
//	return Where{key: key, value: value}
//}
//
//func (w Where) Key() string {
//	return w.key
//}
//
//func (w Where) EQ(item map[string]interface{}) bool {
//	return item[w.key] == w.value[0]
//}
//
//func (w Where) Exec(item map[string]interface{}) bool {
//	return true
//}
