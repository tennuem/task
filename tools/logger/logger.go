package logger

import (
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitloglevel "github.com/go-kit/kit/log/level"
)

func NewLogger(level, format string) kitlog.Logger {
	lvl, err := getLevel(level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init: %s", err)
		os.Exit(1)
	}

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitloglevel.NewFilter(logger, lvl)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	return logger
}

func getLevel(lvl string) (level.Option, error) {
	switch lvl {
	case "emerg":
		return kitloglevel.AllowError(), nil
	case "alert":
		return kitloglevel.AllowError(), nil
	case "crit":
		return kitloglevel.AllowError(), nil
	case "err":
		return kitloglevel.AllowError(), nil
	case "warn":
		return kitloglevel.AllowWarn(), nil
	case "notice":
		return kitloglevel.AllowInfo(), nil
	case "info":
		return kitloglevel.AllowInfo(), nil
	case "debug":
		return kitloglevel.AllowDebug(), nil
	}
	return nil, fmt.Errorf("level %s is incorrect. Level can be (emerg, alert, crit, err, warn, notice, info, debug)", lvl)
}
