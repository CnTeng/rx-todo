package daemon

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/rpc"
	"golang.org/x/exp/jsonrpc2"
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
	registerProjectHandles(s)

	return s.Serve(ctx)
}

func simpleHandle(ctx context.Context, req *jsonrpc2.Request, action func(*client.Client, int64) (any, error)) (any, error) {
	r := &struct {
		ID int64 `json:"id"`
	}{}
	if err := json.Unmarshal(req.Params, &r); err != nil {
		return nil, err
	}

	c, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return action(c, r.ID)
}
