package main

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/tennuem/task/cmd/server/http"
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

	quit := make(chan bool)
	http.Run(cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, mux, quit, log.With(logger, "component", "http server"))
	<-quit
}
