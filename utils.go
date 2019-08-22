package jsonrpc

import (
	"github.com/danhper/structomap"
	"reflect"
	"time"
)

var structSerializer = structomap.New().PickAll()

func getJsonField(m JsonMap, name string, out interface{}) bool {
	var val interface{}
	if val = m[name]; val == nil {
		return false
	}

	ret := false
	switch out := out.(type) {
	case *string:
		if n, ok := val.(string); ok {
			*out = n
			ret = true
		}
	case *float32:
		if n, ok := val.(float64); ok {
			*out = float32(n)
			ret = true
		}
	case *float64:
		if n, ok := val.(float64); ok {
			*out = n
			ret = true
		}
	case *int:
		if n, ok := val.(float64); ok {
			*out = int(n)
			ret = true
		}
	case *uint:
		if n, ok := val.(float64); ok {
			*out = uint(n)
			ret = true
		}
	case *int64:
		if n, ok := val.(float64); ok {
			*out = int64(n)
			ret = true
		}
	case *uint64:
		if n, ok := val.(float64); ok {
			*out = uint64(n)
			ret = true
		}
	case *time.Time:
		if s, ok := val.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				*out = t
				ret = true
			}
		}
	}
	return ret
}

func objectToMap(st interface{}) map[string]interface{} {
	v := reflect.ValueOf(st)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		return structSerializer.Transform(st)
	}
	return nil
}
