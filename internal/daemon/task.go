package daemon

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/CnTeng/rx-todo/internal/rpc"
	"golang.org/x/exp/jsonrpc2"
)

func taskCreateHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	r := &model.TaskCreationRequest{}
	if err := json.Unmarshal(req.Params, &r); err != nil {
		return nil, err
	}

	c, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return c.CreateTask(r)
}

func taskListHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	c, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return c.GetTasks(), nil
}

func taskUpdateHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	r := &client.TaskUpdateRequestWithID{}
	if err := json.Unmarshal(req.Params, &r); err != nil {
		return nil, err
	}

	c, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return c.UpdateTask(r)
}

func taskMoveHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	r := &client.TaskMoveRequestWithID{}
	if err := json.Unmarshal(req.Params, &r); err != nil {
		return nil, err
	}

	c, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return c.MoveTask(r)
}

func taskOpenHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return simpleHandle(ctx, req, func(c *client.Client, id int64) (any, error) {
		return c.OpenTask(id)
	})
}

func taskCloseHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return simpleHandle(ctx, req, func(c *client.Client, id int64) (any, error) {
		return c.CloseTask(id)
	})
}

func taskArchiveHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return simpleHandle(ctx, req, func(c *client.Client, id int64) (any, error) {
		return c.ArchiveTask(id)
	})
}

func taskUnarchiveHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return simpleHandle(ctx, req, func(c *client.Client, id int64) (any, error) {
		return c.UnarchiveTask(id)
	})
}

func taskDeleteHandle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return simpleHandle(ctx, req, func(c *client.Client, id int64) (any, error) {
		return c.DeleteTask(id)
	})
}

func registerTaskHandles(s *rpc.JRPCServer) {
	s.Register("task.create", taskCreateHandle)
	s.Register("task.list", taskListHandle)
	s.Register("task.update", taskUpdateHandle)
	s.Register("task.move", taskMoveHandle)
	s.Register("task.open", taskOpenHandle)
	s.Register("task.close", taskCloseHandle)
	s.Register("task.archive", taskArchiveHandle)
	s.Register("task.unarchive", taskUnarchiveHandle)
	s.Register("task.delete", taskDeleteHandle)
}
