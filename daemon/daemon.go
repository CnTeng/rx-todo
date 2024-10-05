package daemon

import (
	"context"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/rpc"
)

type clientKeyType string

const clientKey clientKeyType = "client"

type Daemon struct {
	address string

	*client.Client
}

func NewDaemon(addr string, client *client.Client) *Daemon {
	return &Daemon{address: addr, Client: client}
}

func (d *Daemon) Serve() error {
	ctx := context.WithValue(context.Background(), clientKey, d.Client)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := rpc.NewServer("unix", d.address)

	registerLabelHandles(s)
	registerTaskHandles(s)

	return s.Serve(ctx)
}
