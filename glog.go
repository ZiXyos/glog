package glog

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

type LogFormat uint8

type config struct {
	writer        io.Writer
	formatter     LogFormat
	options       *slog.HandlerOptions
	handlers      []slog.Handler
	withtimeStamp bool
	reportCaller  bool
	styles        *styleConfig
}

const (
	DefaultFormatter LogFormat = iota
	JSONFormatter
	TextFormatter
)

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

func NewDefault(handlers ...slog.Handler) (*slog.Logger, error) {
	logger, err := New(
		WithLevel(slog.LevelDebug),
		WithJsonFormat(),
		WithTimeStamp(),
		WithReportCaller(),
		WithStyle(
			WithErrorStyle(),
		),
		WithHandler(handlers...),
	)

	return logger, err
}

func createLogger(cfg *config) (*slog.Logger, error) {
	handler := log.New(cfg.writer)
	styles := log.DefaultStyles()

	// Guard against nil styles: only assign if WithStyle was passed
	if cfg.styles != nil {
		styles.Levels = cfg.styles.level
	}

	handler.SetReportTimestamp(cfg.withtimeStamp)
	handler.SetReportCaller(cfg.reportCaller)
	handler.SetStyles(styles)

	// charmbracelet/log.Level values are numerically identical to slog.Level:
	// Debug -4, Info 0, Warn 4, Error 8. Direct cast is safe and intentional.
	level := slog.LevelInfo // default
	if cfg.options.Level != nil {
		level = cfg.options.Level.Level()
	}
	handler.SetLevel(log.Level(level))

	switch cfg.formatter {
	case DefaultFormatter:
		// No-op: charmbracelet's built-in default is already the pretty text formatter
	case JSONFormatter:
		handler.SetFormatter(log.JSONFormatter)
	case TextFormatter:
		handler.SetFormatter(log.TextFormatter)
	default:
		return nil, fmt.Errorf("%s: %d", "unknown format:", cfg.formatter)
	}

	// Build final handler with terminal handler first, then injected handlers
	finalHandler := newMultiHandler(append([]slog.Handler{handler}, cfg.handlers...)...)
	logger := slog.New(finalHandler)

	return logger, nil
}
