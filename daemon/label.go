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

func labelListHandle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.GetLabels(), nil
}

func labelCreateHandle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	r := &model.LabelCreationRequest{}

	if err := json.Unmarshal(req.Params, r); err != nil {
		return nil, err
	}

	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.CreateLabel(r)
}

func labelUpdateHandle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	r := &model.LabelUpdateRequestWithID{}

	if err := json.Unmarshal(req.Params, r); err != nil {
		return nil, err
	}

	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.UpdateLabel(r)
}

func labelDeleteHandle(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	r := &model.LabelDeleteRequestWithID{}

	if err := json.Unmarshal(req.Params, r); err != nil {
		return nil, err
	}

	client, ok := ctx.Value(clientKey).(*client.Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client.DeleteLabel(r.ID)
}

func registerLabelHandles(s *rpc.JRPCServer) {
	s.Register("Label.List", labelListHandle)
	s.Register("Label.Create", labelCreateHandle)
	s.Register("Label.Update", labelUpdateHandle)
}
