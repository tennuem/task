package health

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeGetLivenessEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.GetLiveness(ctx)
		return GetLivenessResponse{err}, nil
	}
}

func makeGetReadinessEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.GetReadiness(ctx)
		return GetReadinessResponse{err}, nil
	}
}

func makeGetVersionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		buildTime, commit, version := s.GetVersion(ctx)
		return GetVersionResponse{buildTime, commit, version}, nil
	}
}

type GetLivenessRequest struct{}

type GetLivenessResponse struct {
	Err error `json:"error,omitempty"`
}

func (r GetLivenessResponse) error() error { return r.Err }

type GetReadinessRequest struct{}

type GetReadinessResponse struct {
	Err error `json:"error,omitempty"`
}

func (r GetReadinessResponse) error() error { return r.Err }

type GetVersionRequest struct{}

type GetVersionResponse struct {
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	Version   string `json:"version"`
}
