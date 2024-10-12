package daemon

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/model"
	"github.com/CnTeng/rx-todo/rpc"
	"golang.org/x/exp/jsonrpc2"
)

func taskListHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.GetTasks(), nil
}

func taskDeleteHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	r := &model.TaskDeleteRequestWithID{}
	if err := json.Unmarshal(req.Params, &r); err != nil {
		return nil, err
	}

	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.DeleteTask(r.ID)
}

func taskCreateHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	r := &model.TaskCreationRequest{}
	if err := json.Unmarshal(req.Params, &r); err != nil {
		return nil, err
	}

	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.CreateTask(r)
}

func registerTaskHandles(s *rpc.JRPCServer) {
	s.Register("Task.List", taskListHandle)
	s.Register("Task.Create", taskCreateHandle)
	s.Register("Task.Delete", taskDeleteHandle)
}
