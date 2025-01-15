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
}
