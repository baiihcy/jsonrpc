package jsonrpc

import (
	"net"
)

type Context struct {
	conn net.Conn
	Req  *Request
	Resp *Response
}

func (ctx *Context) Param(name string) interface{} {
	if ctx.Req == nil {
		return nil
	}
	return ctx.Req.ParamField(name)
}

func (ctx *Context) Param2(name string, out interface{}) bool {
	if ctx.Req == nil {
		return false
	}
	return ctx.Req.ParamField2(name, out)
}

func (ctx *Context) Result(result interface{}) *Context {
	var resultMap ResultMap
	if m, ok := result.(ResultMap); ok {
		resultMap = m
	} else {
		resultMap = objectToMap(result)
	}
	ctx.Resp = &Response{ctx.Req.Id, resultMap, nil}
	return ctx
}

func (ctx *Context) Error(code ErrorCode, message string, data interface{}) *Context {
	var errorData ErrorData
	if m, ok := data.(ErrorData); ok {
		errorData = m
	} else {
		errorData = objectToMap(data)
	}
	ctx.Resp = &Response{ctx.Req.Id, nil, &ResponseError{code, message, errorData}}
	return ctx
}

func (ctx *Context) InternalError(message string, data interface{}) *Context {
	return ctx.Error(ErrorCodeInternal, message, data)
}

func (ctx *Context) ResetResponse() {
	ctx.Resp = nil
}

func (ctx *Context) Respond() bool {
	return ctx.Resp != nil
}