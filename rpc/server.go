package rpc

import (
	"context"
	"errors"

	"golang.org/x/exp/jsonrpc2"
)

type JRPCServer struct {
	network string
	address string
	handles map[string]jsonrpc2.HandlerFunc
}

func NewServer(network, address string) *JRPCServer {
	return &JRPCServer{
		network: network,
		address: address,
		handles: make(map[string]jsonrpc2.HandlerFunc),
	}
}

func (rs *JRPCServer) Handle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	handle, ok := rs.handles[req.Method]
	if ok {
		return handle(ctx, req)
	}

	return nil, errors.New("method not found")
}

func (rs *JRPCServer) Register(method string, handler jsonrpc2.HandlerFunc) {
	rs.handles[method] = handler
}

func (rs *JRPCServer) Serve(ctx context.Context) error {
	l, err := jsonrpc2.NetListener(
		ctx,
		rs.network,
		rs.address,
		jsonrpc2.NetListenOptions{},
	)
	if err != nil {
		return err
	}
	defer l.Close()

	s, err := jsonrpc2.Serve(
		ctx,
		l,
		jsonrpc2.ConnectionOptions{Handler: rs},
	)
	if err != nil {
		return err
	}

	if err := s.Wait(); err != nil {
		return err
	}

	return nil
}
