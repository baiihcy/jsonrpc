package jsonrpc

type ErrorCode int
const (
	// ErrorCodeParse is parse error code.
	ErrorCodeParse ErrorCode = -32700
	// ErrorCodeInvalidRequest is invalid request error code.
	ErrorCodeInvalidRequest ErrorCode = -32600
	// ErrorCodeMethodNotFound is method not found error code.
	ErrorCodeMethodNotFound ErrorCode = -32601
	// ErrorCodeInvalidParams is invalid params error code.
	ErrorCodeInvalidParams ErrorCode = -32602
	// ErrorCodeInternal is internal error code.
	ErrorCodeInternal ErrorCode = -32603
)

const NewLineString = "\r\n"
var NewLineBytes = []byte(NewLineString)
