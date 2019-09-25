package task

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"
)

// NewMetricService returns an instance of an instrumenting Service.
func NewMetricService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &metricService{counter, latency, s}
}

type metricService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *metricService) Create(ctx context.Context, t *Task) (resp *Task, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("handler", "Create", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
		s.requestLatency.With("handler", "Create", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Create(ctx, t)
}

func (s *metricService) GetList(ctx context.Context) (resp []*Task) {
	defer func(begin time.Time) {
		s.requestCount.With("handler", "GetList", "code", "200").Add(1)
		s.requestLatency.With("handler", "GetList", "code", "200").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetList(ctx)
}

func (s *metricService) GetByID(ctx context.Context, ID uuid.UUID) (resp *Task, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("handler", "GetByID", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
		s.requestLatency.With("handler", "GetByID", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetByID(ctx, ID)
}

func (s *metricService) Update(ctx context.Context, t *Task) (resp *Task, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("handler", "Update", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
		s.requestLatency.With("handler", "Update", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Update(ctx, t)
}

func (s *metricService) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("handler", "Delete", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
		s.requestLatency.With("handler", "Delete", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Delete(ctx, ID)
}
