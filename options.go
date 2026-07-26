package glog

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Option provide a function to set configuration
type Option func(*config) error

// WithWriter sets the destination for the terminal handler. Defaults to
// os.Stdout. Useful for redirecting output to a file or capturing it in tests.
func WithWriter(w io.Writer) Option {
	return func(cfg *config) error {
		if w == nil {
			return fmt.Errorf("writer is nil")
		}
		cfg.writer = w

		return nil
	}
}

// WithLevel sets the log level for the logger.
func WithLevel(level slog.Level) Option {
	return func(cfg *config) error {
		cfg.options.Level = level

		return nil
	}
}

// WithFormat sets the log format for the logger.
// Deprecated: WithFormat is deprecated, now use WithTextFormat or WithJsonFormat.
func WithFormat(format LogFormat) Option {
	return func(cfg *config) error {
		cfg.formatter = format

		return nil
	}
}

// WithTextFormat set the formatter to text.
func WithTextFormat() Option {
	return func(c *config) error {
		c.formatter = TextFormatter

		return nil
	}
}

// WithJsonFormat set the formatter to json.
func WithJsonFormat() Option {
	return func(c *config) error {
		c.formatter = JSONFormatter

		return nil
	}
}

func WithTimeStamp() Option {
	return func(cfg *config) error {
		cfg.withtimeStamp = true

		return nil
	}
}

func WithReportCaller() Option {
	return func(cfg *config) error {
		cfg.reportCaller = true

		return nil
	}
}

func WithStyle(styleOpts ...Style) Option {
	return func(cfg *config) error {
		styles, err := newStyle(styleOpts...)
		if err != nil {
			return err
		}

		cfg.styles = styles
		return nil
	}
}

func WithErrorStyle() Style {
	return func(sc *styleConfig) error {
		sc.level[log.ErrorLevel] = lipgloss.NewStyle().
			SetString("ERR ").
			Padding(0, 1, 0).
			Background(lipgloss.Color("204")).
			Foreground(lipgloss.Color("0"))
		return nil
	}
}

// WithHandler appends additional slog.Handler instances to the logger.
// Every log record goes to both the terminal (charmbracelet/log) handler
// and each provided handler, enabling multi-destination logging.
func WithHandler(handlers ...slog.Handler) Option {
	return func(cfg *config) error {
		cfg.handlers = handlers
		return nil
	}
}
