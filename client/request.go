package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/CnTeng/rx-todo/internal/model"
)

type request struct {
	endpoint string
	token    string
	body     any
}

func NewRequest(endpoint, token string) *request {
	return &request{
		endpoint: endpoint,
		token:    token,
	}
}

func (r *request) WithPath(path string) *request {
	if strings.HasPrefix(path, "/") {
		r.endpoint += path
	} else {
		r.endpoint += "/" + path
	}

	return r
}

func (r *request) WithID(id int64) *request {
	return r.WithPath(strconv.FormatInt(id, 10))
}

func (r *request) WithParameter(key string, value *int64) *request {
	if value == nil {
		return r
	}

	if strings.Contains(r.endpoint, "?") {
		r.endpoint += "&" + key + "=" + strconv.FormatInt(*value, 10)
	} else {
		r.endpoint += "?" + key + "=" + strconv.FormatInt(*value, 10)
	}

	return r
}

func (r *request) WithBody(body any) *request {
	r.body = body
	return r
}

func (r *request) Get() (io.ReadCloser, error) {
	return r.Execute(http.MethodGet)
}

func (r *request) Post() (io.ReadCloser, error) {
	return r.Execute(http.MethodPost)
}

func (r *request) Put() (io.ReadCloser, error) {
	return r.Execute(http.MethodPut)
}

func (r *request) Delete() (io.ReadCloser, error) {
	return r.Execute(http.MethodDelete)
}

func (r *request) Execute(method string) (io.ReadCloser, error) {
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
