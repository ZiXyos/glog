# Glog

Just a wrapper around slog and [log](https://github.com/charmbracelet/log)

## Usage
see in [exmaple/main.go](https://github.com/ZiXyos/glog/tree/master/example/main.go)
```go
package main

import (
	"log/slog"

	"github.com/zixyos/glog"
);

func main() {
  logger, err := glog.New(
    glog.WithLevel(slog.LevelDebug),
    glog.WithFormat(glog.TextFormatter),
  );
  if err != nil {
    panic(err)
  }
  logger.Info("Hello World");
}
```
