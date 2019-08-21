package jsonrpc

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net"
	"time"
)

type Server struct {
	opts      ServerOptions
	handlers  map[string]HandlerFunc
	ctx       context.Context
	ctxCancel context.CancelFunc
}

type ServerOptions struct {
	Timeout time.Duration
}

var DefaultServerOptions = ServerOptions{
	Timeout: 10 * time.Second,
}

func NewServer(opts ServerOptions) *Server {
	return &Server{opts: opts, handlers:make(map[string]HandlerFunc)}
}

func (s *Server) Handle(method string, handler HandlerFunc) {
	s.handlers[method] = handler
}

func (s *Server) ListenAndServe(network, address string, ctx context.Context) error {
	ln, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = context.Background()
	}
	s.ctx, s.ctxCancel = context.WithCancel(ctx)

	L:
	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Println("Accept client fail:", err)
		}

		go s.ConnHandler(fd)

		select {
		case <- s.ctx.Done():
			break L
		default:
			time.Sleep(time.Second)
		}
	}
	err = ln.Close()
	return nil
}

func (s *Server) Close() error {
	if s.ctxCancel == nil {
		return errors.New("Server did not serve yet")
	}
	s.ctxCancel()
	return nil
}

func (s *Server) ConnHandler(conn net.Conn) {
	var req Request
	var ctx = &Context{conn, &req, nil}
	reader := bufio.NewReader(conn)

	L:
	for {
		bs, err := reader.ReadBytes('\n')
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Timeout() {
				break L
			}
		}

		ctx.ResetResponse()
		if err := json.Unmarshal(bs, &req); err == nil {
			if handler, ok := s.handlers[req.Method]; ok {
				handler(ctx)
			} else {
				log.Println("Method not found:", ctx.Req.Method)
				ctx.Error(ErrorCodeMethodNotFound, "Method not found:" + req.Method, nil)
			}
		} else {
			log.Println("Unmarshal request fail:", err)
			ctx.Error(ErrorCodeParse, err.Error(), nil)
		}

		if ctx.Respond() {
			if err := s.Respond(ctx); err != nil {
				log.Println("Respond fail", err)
			}
		}


		select {
		case <- s.ctx.Done():
			break L
		default:
			time.Sleep(time.Second)
		}
	}
}

func (s *Server) Respond(ctx *Context) error {
	if ctx.Resp == nil {
		return errors.New("No response")
	}
	bs, err := json.Marshal(ctx.Resp)
	if err != nil {
		return errors.Wrap(err, "Marshal response fail")
	}

	ctx.conn.SetWriteDeadline(time.Now().Add(s.opts.Timeout))

	if _, err := ctx.conn.Write(bs); err != nil {
		return errors.Wrap(err, "Write connection fail")
	}
	if _, err := ctx.conn.Write(NewLineBytes); err != nil {
		return errors.Wrap(err, "Write connection fail")
	}
	log.Println("Server Send:", string(bs))

	return nil
}

