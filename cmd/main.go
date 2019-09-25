package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tennuem/task/configs"
	"github.com/tennuem/task/pkg/repository/inmemory"
	"github.com/tennuem/task/pkg/task"
	"github.com/tennuem/task/tools/logger"
)

func main() {
	var (
		cfg    = configs.NewConfig()
		logger = logger.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
		tasks  = inmemory.NewRepository()
	)

	var taskService task.Service
	taskService = task.NewService(tasks)
	taskService = task.NewLoggingService(log.With(logger, "component", "task"), taskService)
	fieldKeys := []string{"handler", "code"}
	taskService = task.NewMetricService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "esp",
			Subsystem: "task_service",
			Name:      "requests_total",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "esp",
			Subsystem: "task_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		taskService,
	)

	r := mux.NewRouter().StrictSlash(false)
	r.PathPrefix("/task").Handler(task.MakeHTTPHandler(taskService, log.With(logger, "component", "http handler")))

	var g run.Group
	{
		addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			level.Error(logger).Log("component", "HTTP server", "addr", addr, "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("component", "HTTP server", "addr", addr, "msg", "listening...")
			return http.Serve(listener, accessControl(r))
		}, func(error) {
			listener.Close()
		})
	}
	{
		addr := fmt.Sprintf("%s:%d", cfg.Metric.Host, cfg.Metric.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			level.Error(logger).Log("component", "metric server", "addr", addr, "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("component", "metric server", "addr", addr, "msg", "listening...")
			http.Handle("/metrics", promhttp.Handler())
			return http.Serve(listener, http.DefaultServeMux)
		}, func(error) {
			listener.Close()
		})
	}
	{
		addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, 8089)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			level.Error(logger).Log("component", "prefix server", "addr", addr, "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("component", "prefix server", "addr", addr, "msg", "listening...")
			r := mux.NewRouter()
			r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(fmt.Sprintf("%s%s", r.Host, r.RequestURI)))
			})
			return http.Serve(listener, accessControl(r))
		}, func(error) {
			listener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
