package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(request)
		v := svc.Echo(req.S)
		return response{v}, nil
	}
}

func makeUppercaseEndpoint(svc EchoService) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(request)
		v := svc.Uppercase(req.S)
		return response{v}, nil
	}
}

func MakeHTTPHandler(s EchoService, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// 	swagger:route POST /todos echo
	r.Methods("POST").Path("/echo").Handler(httptransport.NewServer(
		makeEchoEndpoint(s),
		decodeRequest,
		encodeResponse,
		httptransport.ServerErrorLogger(logger),
	))

	// 	swagger:route POST /uppercase uppercase
	r.Methods("POST").Path("/uppercase").Handler(httptransport.NewServer(
		makeUppercaseEndpoint(s),
		decodeRequest,
		encodeResponse,
		httptransport.ServerErrorLogger(logger),
	))
	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, HTTPError{code: 400, err: err}
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type request struct {
	S string `json:"s"`
}

type response struct {
	S string `json:"s"`
}
