package main

import (
	"log/slog"

	"github.com/zixyos/glog"
)

func main() {
  logger, err := glog.New(
    glog.WithLevel(slog.LevelDebug),
    glog.WithFormat(glog.TextFormatter),
  );
  if err != nil {
    panic(err)
  }

  logger.Info("NOT SO FAST BOY")
}
