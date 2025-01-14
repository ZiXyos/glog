package glog

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

type LogFormat uint8;

type config struct {
  writer    io.Writer
  formatter LogFormat
  options   *slog.HandlerOptions
}

const (
  DefaultFormatter LogFormat = iota;
  JSONFormatter;
  TextFormatter;
);

func New(opts ...Option) (*slog.Logger, error) {
  cfg := config{
    writer: os.Stdout,
    options: &slog.HandlerOptions{
      Level: slog.LevelInfo,
      ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					return slog.Attr{Key: "Message", Value: a.Value}
				}
				return a
			},
    },
  }

  for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

  return createLogger(&cfg)
}

func createLogger(cfg *config) (*slog.Logger, error) {
  handler := log.New(cfg.writer);
  logger := slog.New(handler);

  switch cfg.formatter {
  case JSONFormatter:
    handler.SetFormatter(log.JSONFormatter)
  case TextFormatter:
    handler.SetFormatter(log.TextFormatter)
  default:
    return nil, fmt.Errorf("%s: %d", "unknown format:", cfg.formatter)
}

  return logger, nil;
}

