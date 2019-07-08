package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func Run(host string, port int, handler http.Handler, quit chan bool, logger log.Logger) {
	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr:    addr,
		Handler: accessControl(handler),
	}

	go func() {
		level.Info(logger).Log("address", addr, "msg", "listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			level.Error(logger).Log("err", err.Error())
		}
		time.Sleep(time.Second * 1)
		level.Info(logger).Log("msg", "gracefull stoped")
		quit <- true
	}()

	go shutDown(context.TODO(), srv, logger)
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

func shutDown(ctx context.Context, srv *http.Server, logger log.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	level.Info(logger).Log("msg", "shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		level.Error(logger).Log("err", err.Error())
	}

	level.Info(logger).Log("msg", "Shutting down done")
}
