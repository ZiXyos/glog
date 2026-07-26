package glog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

// TestNewZeroOptionsSucceeds tests that New() with zero options does not return
// "unknown format: 0" error. Previously the zero-value DefaultFormatter fell through
// to the default: error branch.
func TestNewZeroOptionsSucceeds(t *testing.T) {
	buf := &bytes.Buffer{}
	logger, err := New(WithWriter(buf))
	if err != nil {
		t.Fatalf("New() with zero options returned error: %v", err)
	}
	if logger == nil {
		t.Error("New() returned nil logger")
	}
}

// TestNoPanicWithoutWithStyle tests that New() does not nil-deref cfg.styles
// when WithStyle is not provided. Previously this caused a panic.
func TestNoPanicWithoutWithStyle(t *testing.T) {
	buf := &bytes.Buffer{}
	logger, err := New(WithWriter(buf), WithLevel(slog.LevelDebug))
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// Log a record to verify it reaches the buffer without panic
	logger.Debug("test message")
	if buf.Len() == 0 {
		t.Error("expected logged record to reach buffer")
	}
}

// TestWithLevelFiltersInfo tests that WithLevel actually filters log records.
// With WithLevel(slog.LevelWarn), Info records must NOT appear but Warn records must.
func TestWithLevelFiltersInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	logger, err := New(
		WithWriter(buf),
		WithLevel(slog.LevelWarn),
		WithTextFormat(),
	)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	// Log an Info record - should NOT appear
	logger.Info("info message")
	infoContent := buf.String()
	if strings.Contains(infoContent, "info message") {
		t.Error("Info record should not appear with level=Warn")
	}

	// Clear buffer and log a Warn record - should appear
	buf.Reset()
	logger.Warn("warn message")
	warnContent := buf.String()
	if !strings.Contains(warnContent, "warn message") {
		t.Error("Warn record should appear with level=Warn")
	}
}

// TestWithLevelDebugLetsThroughDebug tests that WithLevel(slog.LevelDebug) allows
// Debug records through. charm's default is Info, so this proves the level is applied.
func TestWithLevelDebugLetsThroughDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger, err := New(
		WithWriter(buf),
		WithLevel(slog.LevelDebug),
		WithTextFormat(),
	)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	logger.Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug record should appear with level=Debug")
	}
}

// TestWithWriterNilReturnsError tests that WithWriter(nil) returns an error,
// not a panic.
func TestWithWriterNilReturnsError(t *testing.T) {
	_, err := New(WithWriter(nil))
	if err == nil {
		t.Fatal("WithWriter(nil) should return an error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil, got: %v", err)
	}
}

// TestWithHandlerFanoutReachesAllLegs tests that WithHandler fanout reaches both
// the terminal handler and the additional handler. Log one record and assert
// the message appears in BOTH buffers.
func TestWithHandlerFanoutReachesAllLegs(t *testing.T) {
	terminalBuf := &bytes.Buffer{}
	jsonBuf := &bytes.Buffer{}

	jsonHandler := slog.NewJSONHandler(jsonBuf, nil)
	logger, err := New(
		WithWriter(terminalBuf),
		WithHandler(jsonHandler),
	)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	logger.Info("fanout test message")

	// Check terminal buffer
	if !strings.Contains(terminalBuf.String(), "fanout test message") {
		t.Error("message should appear in terminal buffer")
	}

	// Check JSON handler buffer
	if !strings.Contains(jsonBuf.String(), "fanout test message") {
		t.Error("message should appear in JSON handler buffer")
	}
}

// TestUnknownFormatterErrors tests that New(WithFormat(LogFormat(99))) returns
// an error, proving the default: branch was kept.
func TestUnknownFormatterErrors(t *testing.T) {
	buf := &bytes.Buffer{}
	_, err := New(WithWriter(buf), WithFormat(LogFormat(99)))
	if err == nil {
		t.Fatal("New with unknown format should return an error")
	}
	if !strings.Contains(err.Error(), "unknown format") {
		t.Errorf("error should mention unknown format, got: %v", err)
	}
}

// TestNewDefaultSucceeds tests that NewDefault() succeeds and writes to a
// provided handler. Since NewDefault writes its terminal leg to os.Stdout
// (it takes no writer), we assert only on the handler's buffer.
func TestNewDefaultSucceeds(t *testing.T) {
	jsonBuf := &bytes.Buffer{}
	jsonHandler := slog.NewJSONHandler(jsonBuf, nil)

	logger, err := NewDefault(jsonHandler)
	if err != nil {
		t.Fatalf("NewDefault() returned error: %v", err)
	}
	if logger == nil {
		t.Fatal("NewDefault() returned nil logger")
	}

	logger.Info("default logger message")

	// Assert the message reached the JSON handler
	if !strings.Contains(jsonBuf.String(), "default logger message") {
		t.Error("message should appear in JSON handler buffer")
	}
}
