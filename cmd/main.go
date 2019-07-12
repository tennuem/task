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
	"github.com/gorilla/mux"
	"github.com/oklog/run"
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

	mux := mux.NewRouter().StrictSlash(false)
	mux.PathPrefix("/task").Handler(task.MakeHTTPHandler(taskService, log.With(logger, "component", "http handler")))

	var g run.Group
	{
		addr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			level.Error(logger).Log("transport", "HTTP", "addr", addr, "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Info(logger).Log("transport", "HTTP", "addr", addr, "msg", "listening")
			return http.Serve(listener, accessControl(mux))
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
