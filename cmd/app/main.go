package main

import (
	"context"
	"github.com/tennuem/task/internal/server"
)

func main() {
	svr := server.NewServer()
	svr.Run(context.Background())
}
