package rpc

import (
	"context"
	"net"
	"time"

	"golang.org/x/exp/jsonrpc2"
)

type JRPCClient struct {
	network string
	address string
	timeout time.Duration
}

func NewClient(network, address string, timeout time.Duration) *JRPCClient {
	return &JRPCClient{
		network: network,
		address: address,
		timeout: timeout,
	}
}

func (rc *JRPCClient) Call(method string, params, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), rc.timeout)
	defer cancel()

	conn, err := jsonrpc2.Dial(
		ctx,
		jsonrpc2.NetDialer(rc.network, rc.address, net.Dialer{}),
		jsonrpc2.ConnectionOptions{},
	)
	if err != nil {
		return err
	}

	call := conn.Call(ctx, method, params)

	if err := call.Await(ctx, &result); err != nil {
		return err
	}

	return nil
}
