package glog

import (
	"log/slog"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Option provide a function to set configuration
type Option func(*config) error;

// WithLevel sets the log level for the logger.
func WithLevel(level slog.Level) Option {
	return func(cfg *config) error {
		cfg.options.Level = level;

		return nil
	}
}

// WithFormat sets the log format for the logger.
func WithFormat(format LogFormat) Option {
	return func(cfg *config) error {
		cfg.formatter = format;

		return nil
	}
}

func WithTimeStamp() Option {
  return func(cfg *config) error {
    cfg.withtimeStamp = true; 

    return nil
  }
}

func WithReportCaller() Option {
  return func(cfg *config) error {
    cfg.reportCaller = true;

    return nil
  }
}

func WithStyle(styleOpts ...Style) Option {
  return func(cfg *config) error {
    styles, err := newStyle(styleOpts...)
    if err != nil {
      return err
    }

    cfg.styles = styles;
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
