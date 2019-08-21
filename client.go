package jsonrpc

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	opts ClientOptions
	lock sync.Mutex
	idCounter int
}

type ClientOptions struct {
	Timeout time.Duration
}

var DefaultClientOptions = ClientOptions{
	Timeout: 10 * time.Second,
}

func NewClient(opts ClientOptions) *Client {
	return &Client{opts: opts}
}

func (c *Client) Dial(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Request(method string, params interface{}) (*Response, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var paramMap ParamMap
	if m, ok := params.(ParamMap); ok {
		paramMap = m
	} else {
		paramMap = objectToMap(params)
	}

	c.idCounter++
	req := &Request{c.idCounter, method, paramMap}
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(c.opts.Timeout)
	c.conn.SetReadDeadline(deadline)
	c.conn.SetWriteDeadline(deadline)

	if _, err = c.conn.Write(bs); err != nil {
		return nil, errors.Wrap(err, "Write connection fail")
	}
	if _, err = c.conn.Write(NewLineBytes); err != nil {
		return nil, errors.Wrap(err, "Write connection fail")
	}
	log.Println("Client send:", string(bs))

	reader := bufio.NewReader(c.conn)
	bs, err = reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var resp Response
	if err = json.Unmarshal(bs, &resp); err != nil {
		return nil, err
	}
	return &resp, err
}
