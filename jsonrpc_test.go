package jsonrpc

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJsonrpc(t *testing.T) {
	type TestEntity struct {
		A int
		B string
		C time.Time
	}

	now := time.Now()
	s := NewServer(DefaultServerOptions)
	defer s.Close()
	go s.ListenAndServe("tcp", ":55667", nil)
	s.Handle("testResult", func(ctx *Context) error {
		var paramA float64
		var paramB string
		var paramC time.Time
		assert.False(t, ctx.Req.ParamField2("A", &paramB))
		assert.True(t, ctx.Req.ParamField2("A", &paramA), "Param A not found or not a number")
		assert.True(t, ctx.Req.ParamField2("B", &paramB), "Param B not found or not a string")
		assert.True(t, ctx.Req.ParamField2("C", &paramC), "Param C not found or not a time")
		assert.Equal(t, 11., paramA)
		assert.Equal(t, "params", paramB)

		t.Log(ctx.Req)
		ctx.Result(&TestEntity{-11, "result", now})
		return nil
	})
	s.Handle("testError", func(ctx *Context) error {
		ctx.InternalError("test error", 999)
		return nil
	})
	s.Handle("testInternalError", func(ctx *Context) error {
		return errors.New("internal error")
	})

	c := NewClient(DefaultClientOptions)
	if err := c.Dial("tcp", "127.0.0.1:55667"); err != nil {
		t.Fatal("Connect to server fail:", err)
	}

	var resultA int
	var resultB string
	var resultC time.Time
	resp, err := c.Request("testResult", TestEntity{11, "params", now})
	if err != nil {
		t.Error("Request testResult fail:", err)
	}
	assert.Nil(t, resp.Error)
	assert.False(t, resp.ResultField2("A", &resultB))
	assert.True(t, resp.ResultField2("A", &resultA), "Result field A not found or not a number")
	assert.True(t, resp.ResultField2("B", &resultB), "Result field B not found or not a string")
	assert.True(t, resp.ResultField2("C", &resultC), "Result field C not found or not a time")
	assert.Equal(t, -11, resultA)
	assert.Equal(t, "result", resultB)
	t.Logf("%+v\n", resp)

	resp, err = c.Request("testError", nil)
	if err != nil {
		t.Error("Request testError fail:", err)
	}
	assert.NotNil(t, resp.Error)

	resp, err = c.Request("testInternalError", nil)
	if err != nil {
		t.Error("Request testInternalError fail:", err)
	}
	assert.NotNil(t, resp.Error)
	assert.Equal(t, "internal error", resp.Error.Message)
}
