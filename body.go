package gocrud_redis

import "github.com/kordar/gocrud"

type RedisSearchBody struct {
	*gocrud.SearchBody
}

func NewRedisSearchBody(body gocrud.SearchBody) RedisSearchBody {
	return RedisSearchBody{SearchBody: &body}
}

func (search *RedisSearchBody) RedisQuery(db interface{}, parallel map[string]string) Table {
	out := search.Query(db, parallel)
	if t, ok := out.(Table); ok {
		return t
	}
	if t, ok := db.(Table); ok {
		return t
	}
	return nil
}

func (search *RedisSearchBody) RedisPaginate(db interface{}, parallel map[string]string) (Table, error) {
	tx, err := search.Paginate(db, parallel)
	if tx == nil {
		return nil, err
	}
	if t, ok := tx.(Table); ok {
		return t, err
	}
	if t, ok := db.(Table); ok {
		return t, err
	}
	return nil, err
}

func (search *RedisSearchBody) RedisQueryCustom(f func(search *gocrud.SearchBody) interface{}) Table {
	out := search.QueryCustom(f)
	if t, ok := out.(Table); ok {
		return t
	}
	return nil
}

type RedisEditorBody struct {
	*gocrud.EditorBody
}

func NewRedisEditorBody(body gocrud.EditorBody) RedisEditorBody {
	return RedisEditorBody{EditorBody: &body}
}

func (form *RedisEditorBody) RedisQuery(db interface{}, parallel map[string]string) Table {
	out := form.Query(db, parallel)
	if t, ok := out.(Table); ok {
		return t
	}
	if t, ok := db.(Table); ok {
		return t
	}
	return nil
}

func (form *RedisEditorBody) RedisQuerySafe(db interface{}, parallel map[string]string) (Table, error) {
	tx, err := form.QuerySafe(db, parallel)
	if tx == nil {
		return nil, err
	}
	if t, ok := tx.(Table); ok {
		return t, err
	}
	return nil, err
}

func (form *RedisEditorBody) RedisQueryCustom(f func(body *gocrud.EditorBody) interface{}) Table {
	out := form.QueryCustom(f)
	if t, ok := out.(Table); ok {
		return t
	}
	return nil
}

type RedisRemoveBody struct {
	*gocrud.RemoveBody
}

func NewRedisRemoveBody(body gocrud.RemoveBody) RedisRemoveBody {
	return RedisRemoveBody{RemoveBody: &body}
}

func (remove *RedisRemoveBody) RedisQuerySafe(db interface{}, parallel map[string]string) (Table, error) {
	tx, err := remove.QuerySafe(db, parallel)
	if tx == nil {
		return nil, err
	}
	if t, ok := tx.(Table); ok {
		return t, err
	}
	return nil, err
}

func (remove *RedisRemoveBody) RedisQuery(db interface{}, parallel map[string]string) Table {
	out := remove.Query(db, parallel)
	if t, ok := out.(Table); ok {
		return t
	}
	if t, ok := db.(Table); ok {
		return t
	}
	return nil
}

func (remove *RedisRemoveBody) RedisQueryCustom(f func(body *gocrud.RemoveBody) interface{}) Table {
	out := remove.QueryCustom(f)
	if t, ok := out.(Table); ok {
		return t
	}
	return nil
}

type RedisFormBody struct {
	*gocrud.FormBody
}

func NewRedisFormBody(body gocrud.FormBody) RedisFormBody {
	return RedisFormBody{FormBody: &body}
}

func (form *RedisFormBody) RedisQuery(db interface{}, parallel map[string]string) Table {
	out := form.Query(db, parallel)
	if t, ok := out.(Table); ok {
		return t
	}
	if t, ok := db.(Table); ok {
		return t
	}
	return nil
}

func (form *RedisFormBody) RedisQuerySafe(db interface{}, parallel map[string]string) (Table, error) {
	tx, err := form.QuerySafe(db, parallel)
	if tx == nil {
		return nil, err
	}
	if t, ok := tx.(Table); ok {
		return t, err
	}
	return nil, err
}

func (form *RedisFormBody) RedisQueryCustom(f func(body *gocrud.FormBody) interface{}) Table {
	out := form.QueryCustom(f)
	if t, ok := out.(Table); ok {
		return t
	}
	return nil
}
