package jsonrpc


type Request struct {
	Id int
	Method string
	Params ParamMap `json:",omitempty"`
}

func (req *Request) ParamField(name string) interface{} {
	return req.Params[name]
}

func (req *Request) ParamField2(name string, out interface{}) bool {
	return getJsonField(req.Params, name, out)
}
