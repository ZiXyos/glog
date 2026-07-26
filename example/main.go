package main

import (
	"log/slog"
	"os"

	"github.com/zixyos/glog"
)

func main() {
	logger, err := glog.New(
		glog.WithLevel(slog.LevelDebug),
		glog.WithFormat(glog.TextFormatter),
		glog.WithReportCaller(),
		glog.WithTimeStamp(),
		glog.WithStyle(
			glog.WithErrorStyle(),
		),
	)
	if err != nil {
		panic(err)
	}

	logger.Info("NOT SO FAST BOY")
	logger.Error("THIS IS BAd")

	defLogger, err := glog.NewDefault()
	if err != nil {
		panic(err)
	}

	defLogger.Warn("nothing to see here", "sender", "your dear maintainer")

	// Fanout example: send logs to both terminal and a JSON handler.
	// The second handler can be any slog.Handler, e.g. an OpenTelemetry one.
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	fanoutLogger, err := glog.New(
		glog.WithLevel(slog.LevelDebug),
		glog.WithHandler(jsonHandler),
	)
	if err != nil {
		panic(err)
	}

	fanoutLogger.Info("This appears in both terminal and JSON output")
}
