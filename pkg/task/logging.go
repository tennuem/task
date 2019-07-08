package task

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

type loggingService struct {
	logger log.Logger
	Service
}

func (s *loggingService) Create(ctx context.Context, t *Task) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(err),
			"method", "Create",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Create(ctx, t)
}

func (s *loggingService) GetList(ctx context.Context) []*Task {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(nil),
			"method", "GetList",
			"took", time.Since(begin),
			"err", nil,
		)
	}(time.Now())
	return s.Service.GetList(ctx)
}

func (s *loggingService) GetByID(ctx context.Context, ID int) (t *Task, err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(err),
			"method", "GetByID",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetByID(ctx, ID)
}

func (s *loggingService) Update(ctx context.Context, t *Task) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(err),
			"method", "Update",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, t)
}

func (s *loggingService) Delete(ctx context.Context, ID int) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(err),
			"method", "Delete",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, ID)
}
