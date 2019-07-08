package task

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		err := s.Create(ctx, req.Task)
		return CreateResponse{Err: err}, nil
	}
}

func makeGetListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp := s.GetList(ctx)
		return GetListResponse{Tasks: resp}, nil
	}
}

func makeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		resp, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{resp, err}, nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := s.Update(ctx, req.Task)
		return UpdateResponse{err}, nil
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.ID)
		return DeleteResponse{err}, nil
	}
}

type CreateRequest struct {
	Task *Task `json:"task,omitempty"`
}

type CreateResponse struct {
	Err error `json:"error,omitempty"`
}

func (r CreateResponse) error() error { return r.Err }

type GetListRequest struct{}

type GetListResponse struct {
	Tasks []*Task `json:"tasks,omitempty"`
	Err   error   `json:"error,omitempty"`
}

func (r GetListResponse) error() error { return r.Err }

type GetByIDRequest struct {
	ID int `json:"id,omitempty"`
}

type GetByIDResponse struct {
	Task *Task `json:"task,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r GetByIDResponse) error() error { return r.Err }

type UpdateRequest struct {
	Task *Task `json:"task,omitempty"`
}

type UpdateResponse struct {
	Err error `json:"error,omitempty"`
}

func (r UpdateResponse) error() error { return r.Err }

type DeleteRequest struct {
	ID int `json:"id,omitempty"`
}

type DeleteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteResponse) error() error { return r.Err }
