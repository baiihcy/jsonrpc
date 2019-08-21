package jsonrpc


type HandlerFunc func (ctx *Context) error
type ParamMap = map[string]interface{}
type ResultMap = map[string]interface{}
type ErrorData = map[string]interface{}
type JsonMap = map[string]interface{}