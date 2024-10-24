package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/CnTeng/rx-todo/model"
)

type request struct {
	endpoint string
	token    string
	body     any
}

func newRequest(endpoint, token string, body any) *request {
	return &request{
		endpoint: endpoint,
		token:    token,
		body:     body,
	}
}

func (r *request) withPath(path string) *request {
	if strings.HasPrefix(path, "/") {
		r.endpoint += path
	} else {
		r.endpoint += "/" + path
	}

	return r
}

func (r *request) withID(id int64) *request {
	return r.withPath(strconv.FormatInt(id, 10))
}

func (r *request) get() (io.ReadCloser, error) {
	return r.execute(http.MethodGet)
}

func (r *request) post() (io.ReadCloser, error) {
	return r.execute(http.MethodPost)
}

func (r *request) put() (io.ReadCloser, error) {
	return r.execute(http.MethodPut)
}

func (r *request) delete() (io.ReadCloser, error) {
	return r.execute(http.MethodDelete)
}

func (r *request) execute(method string) (io.ReadCloser, error) {
	b, err := r.marshalBody()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, r.endpoint, b)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+r.token)

	c := http.Client{}
	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNoContent {
		response.Body.Close()
		return nil, nil
	}

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusCreated {
		var errResp model.ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		response.Body.Close()
		return nil, errResp.Error()
	}

	return response.Body, nil
}

func (r *request) marshalBody() (io.ReadCloser, error) {
	if r.body == nil {
		return http.NoBody, nil
	}

	if rc, ok := r.body.(io.ReadCloser); ok {
		return rc, nil
	}

	if err := model.Validate(r.body); err != nil {
		return nil, err
	}

	b, err := json.Marshal(r.body)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewBuffer(b)), nil
}
