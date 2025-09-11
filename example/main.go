package main

import (
	"log/slog"

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
  );
  if err != nil {
    panic(err)
  }

  logger.Info("NOT SO FAST BOY")
  logger.Error("THIS IS BAd")

	defLogger, err := glog.NewDefault()

	defLogger.Warn("nothing to see here", "sender", "your dear maintainer")
}
