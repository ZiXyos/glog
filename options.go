package glog

import "log/slog"

// Option provide a function to set configuration
type Option func(*config) error


// WithLevel sets the log level for the logger.
func WithLevel(level slog.Level) Option {
	return func(cfg *config) error {
		cfg.options.Level = level

		return nil
	}
}

// WithFormat sets the log format for the logger.
func WithFormat(format LogFormat) Option {
	return func(cfg *config) error {
		cfg.formatter = format

		return nil
	}
}
