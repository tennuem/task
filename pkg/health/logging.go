package health

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

type loggingService struct {
	logger log.Logger
	Service
}

func (s *loggingService) GetLiveness(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(err),
			"method", "GetLiveness",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetLiveness(ctx)
}

func (s *loggingService) GetReadiness(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(err),
			"method", "GetReadiness",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetReadiness(ctx)
}

func (s *loggingService) GetVersion(ctx context.Context) (buildTime, commit, version string) {
	defer func(begin time.Time) {
		level.Info(s.logger).Log(
			"code", getHTTPStatusCode(nil),
			"method", "GetVersion",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetVersion(ctx)
}
