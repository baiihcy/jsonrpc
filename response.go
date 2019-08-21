package jsonrpc

type Response struct {
	Id int
	Result ResultMap `json:",omitempty"`
	Error *ResponseError `json:",omitempty"`
}

type ResponseError struct {
	Code ErrorCode
	Message string
	Data ErrorData `json:",omitempty"`
}

func (resp *Response) ResultField(name string) interface{} {
	return resp.Result[name]
}

func (resp *Response) ResultField2(name string, out interface{}) bool {
	return getJsonField(resp.Result, name, out)
}

func (resp *Response) ErrorField(name string) interface{} {
	if resp.Error == nil {
		return nil
	}
	return resp.Error.Data[name]
}

func (resp *Response) ErrorField2(name string, out interface{}) bool {
	if resp.Error == nil {
		return false
	}
	return getJsonField(resp.Error.Data, name, out)
}
