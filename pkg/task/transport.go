package task

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var errBadRoute = errors.New("bad route")

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
	}

	r.Methods("POST").Path("/task").Handler(httptransport.NewServer(
		makeCreateEndpoint(s),
		decodeHTTPCreateRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/task").Handler(httptransport.NewServer(
		makeGetListEndpoint(s),
		decodeHTTPGetListRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/task/{id}").Handler(httptransport.NewServer(
		makeGetByIDEndpoint(s),
		decodeHTTPGetByIDRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/task/{id}").Handler(httptransport.NewServer(
		makeUpdateEndpoint(s),
		decodeHTTPUpdateRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/task/{id}").Handler(httptransport.NewServer(
		makeDeleteEndpoint(s),
		decodeHTTPDeleteRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeHTTPGetListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetByIDRequest{}, nil
}

func decodeHTTPGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	s, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil, errBadRoute
	}
	return GetByIDRequest{ID: int(id)}, nil
}

func decodeHTTPUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeHTTPDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	s, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil, errBadRoute
	}
	return DeleteRequest{ID: int(id)}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(getHTTPStatusCode(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func getHTTPStatusCode(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
