package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
	}

	r.Methods("GET").Path("/liveness").Handler(httptransport.NewServer(
		makeGetLivenessEndpoint(s),
		decodeGetLivenessRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/readiness").Handler(httptransport.NewServer(
		makeGetReadinessEndpoint(s),
		decodeGetReadinessRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/version").Handler(httptransport.NewServer(
		makeGetVersionEndpoint(s),
		decodeGetVersionRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeGetLivenessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetLivenessRequest{}, nil
}

func decodeGetReadinessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetReadinessRequest{}, nil
}

func decodeGetVersionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetVersionRequest{}, nil
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
