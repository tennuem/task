package server

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/tennuem/task/configs"
	"github.com/tennuem/task/pkg/repository/inmemory"
	"github.com/tennuem/task/pkg/task"
	"github.com/tennuem/task/tools/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Run(ctx context.Context) error
}

func NewServer() Server {
	cfg := configs.NewConfig()
	logger := logger.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
	tasks := inmemory.NewRepository()

	var taskService task.Service
	taskService = task.NewService(tasks)
	taskService = task.NewLoggingService(log.With(logger, "component", "task"), taskService)

	r := mux.NewRouter().StrictSlash(false)
	r.PathPrefix("/task").Handler(task.MakeHTTPHandler(taskService, log.With(logger, "component", "http handler")))

	return &server{
		cfg:     cfg,
		logger:  log.With(logger, "component", "server"),
		handler: r,
	}
}

type server struct {
	g       run.Group
	cfg     *configs.Config
	logger  log.Logger
	handler http.Handler
}

func (s *server) Run(ctx context.Context) error {
	s.shutdown(ctx)
	if err := s.runHTTP(); err != nil {
		return level.Error(s.logger).Log("err", errors.Wrap(err, "run http server"))
	}
	s.interruptSignal()
	return s.logger.Log("exit", s.g.Run())
}

func (s *server) runHTTP() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.HTTP.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.g.Add(func() error {
		level.Info(s.logger).Log("msg", fmt.Sprintf("run http server on %s", addr))
		return http.Serve(listener, accessControl(s.handler))
	}, func(error) {
		listener.Close()
	})
	return nil
}

func (s *server) interruptSignal() {
	ch := make(chan struct{})
	s.g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-ch:
			return nil
		}
	}, func(error) {
		close(ch)
	})
}

func (s *server) shutdown(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	s.g.Add(func() error {
		<-ctx.Done()
		return errors.New("shutdown done")
	}, func(error) {
		cancel()
	})
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
