package jsonrpc

import (
	conv "github.com/cstockton/go-conv"
	"github.com/danhper/structomap"
	"log"
	"reflect"
)

var structSerializer = structomap.New().PickAll()

func getJsonField(m JsonMap, name string, out interface{}) bool {
	var val interface{}
	if val = m[name]; val == nil {
		return false
	}

	if err := conv.Infer(out, val); err != nil {
		log.Println(err)
		return false
	}
	return true
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
