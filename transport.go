package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

// HTTPError is an error type used to customize the status code returned by the HTTP API
type HTTPError struct {
	code int
	err  error
}

// StatusCode returns the status code of the HTTP response.
// It implements the StatusCoder interface
func (e HTTPError) StatusCode() int {
	return e.code
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("status %d, err: %v", e.code, e.err)
}

func makeEchoEndpoint(svc EchoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(echoRequest)
		v := svc.Echo(req.S)
		return echoResponse{v}, nil
	}
}

func decodeEchoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request echoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, HTTPError{code: 400, err: err}
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type echoRequest struct {
	S string `json:"s"`
}

type echoResponse struct {
	S string `json:"s"`
}
