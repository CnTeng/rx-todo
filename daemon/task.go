package daemon

import (
	"context"
	"errors"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/rpc"
	"golang.org/x/exp/jsonrpc2"
)

func taskListHandle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.GetTasks(), nil
}

func registerTaskHandles(s *rpc.JRPCServer) {
	s.Register("Task.List", taskListHandle)
}
