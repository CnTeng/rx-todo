package daemon

import (
	"context"
	"errors"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/internal/rpc"
	"golang.org/x/exp/jsonrpc2"
)

func projectListHandle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.GetProjects(), nil
}

func registerProjectHandles(s *rpc.JRPCServer) {
	s.Register("project.list", projectListHandle)
}
