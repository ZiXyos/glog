package glog

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

// recordingHandler is a fake handler that records the records and attrs it received.
type recordingHandler struct {
	records        []slog.Record
	attrs          []slog.Attr
	enabledResult  bool
	handleErr      error
	withAttrsAttrs []slog.Attr
	withGroupName  string
}

func (rh *recordingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return rh.enabledResult
}

func (rh *recordingHandler) Handle(ctx context.Context, r slog.Record) error {
	rh.records = append(rh.records, r)
	return rh.handleErr
}

func (rh *recordingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newRh := &recordingHandler{
		records:       rh.records,
		attrs:         append([]slog.Attr{}, attrs...),
		enabledResult: rh.enabledResult,
		handleErr:     rh.handleErr,
	}
	return newRh
}

func (rh *recordingHandler) WithGroup(name string) slog.Handler {
	newRh := &recordingHandler{
		records:       rh.records,
		attrs:         rh.attrs,
		enabledResult: rh.enabledResult,
		handleErr:     rh.handleErr,
		withGroupName: name,
	}
	return newRh
}

func TestNilChildrenFiltered(t *testing.T) {
	// Passing a nil handler alongside a real one must not panic on Handle
	realHandler := &recordingHandler{enabledResult: true}
	h := newMultiHandler(nil, realHandler, nil)

	record := slog.Record{}
	err := h.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	if len(realHandler.records) != 1 {
		t.Errorf("expected 1 record, got %d", len(realHandler.records))
	}
}

func TestRecordDispatchedToAllChildren(t *testing.T) {
	// A record dispatched through the fanout reaches ALL children
	h1 := &recordingHandler{enabledResult: true}
	h2 := &recordingHandler{enabledResult: true}
	h3 := &recordingHandler{enabledResult: true}

	fanout := newMultiHandler(h1, h2, h3)

	record := slog.Record{}
	err := fanout.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	if len(h1.records) != 1 {
		t.Errorf("h1: expected 1 record, got %d", len(h1.records))
	}
	if len(h2.records) != 1 {
		t.Errorf("h2: expected 1 record, got %d", len(h2.records))
	}
	if len(h3.records) != 1 {
		t.Errorf("h3: expected 1 record, got %d", len(h3.records))
	}
}

func TestErrorIsolation(t *testing.T) {
	// Error isolation: if the first child returns an error, the second child
	// still receives the record, and the returned error wraps the first child's error
	errFirst := errors.New("first error")
	h1 := &recordingHandler{enabledResult: true, handleErr: errFirst}
	h2 := &recordingHandler{enabledResult: true}

	fanout := newMultiHandler(h1, h2)

	record := slog.Record{}
	err := fanout.Handle(context.Background(), record)

	// Both handlers should have received the record
	if len(h1.records) != 1 {
		t.Errorf("h1: expected 1 record, got %d", len(h1.records))
	}
	if len(h2.records) != 1 {
		t.Errorf("h2: expected 1 record, got %d", len(h2.records))
	}

	// Error should contain the first child's error
	if !errors.Is(err, errFirst) {
		t.Errorf("expected error to wrap %v, got %v", errFirst, err)
	}
}

func TestWithAttrsDoesNotAlias(t *testing.T) {
	// WithAttrs does not alias — calling WithAttrs on the fanout returns a new handler
	// whose children differ from the original's
	h1 := &recordingHandler{enabledResult: true}
	h2 := &recordingHandler{enabledResult: true}

	fanout := newMultiHandler(h1, h2)
	attrs := []slog.Attr{slog.String("key", "value")}
	newFanout := fanout.WithAttrs(attrs)

	// Original fanout's children must not have received the attrs
	if len(h1.attrs) != 0 {
		t.Errorf("original h1: expected 0 attrs, got %d", len(h1.attrs))
	}
	if len(h2.attrs) != 0 {
		t.Errorf("original h2: expected 0 attrs, got %d", len(h2.attrs))
	}

	// New fanout should have new children with attrs
	if h, ok := newFanout.(*multiHandler); ok {
		if len(h.children) != 2 {
			t.Errorf("new fanout: expected 2 children, got %d", len(h.children))
		}
		// Verify new children are different instances
		if h.children[0] == h1 {
			t.Error("new fanout's first child should not be the original h1")
		}
		if h.children[1] == h2 {
			t.Error("new fanout's second child should not be the original h2")
		}
	}
}

func TestEnabledIsOR(t *testing.T) {
	// Enabled is an OR: one child disabled + one enabled → fanout Enabled is true
	hDisabled := &recordingHandler{enabledResult: false}
	hEnabled := &recordingHandler{enabledResult: true}

	fanout := newMultiHandler(hDisabled, hEnabled)

	if !fanout.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected fanout.Enabled to be true when at least one child is enabled")
	}

	// All disabled → false
	fanoutAllDisabled := newMultiHandler(
		&recordingHandler{enabledResult: false},
		&recordingHandler{enabledResult: false},
	)

	if fanoutAllDisabled.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected fanout.Enabled to be false when all children are disabled")
	}

	// Zero children → false
	fanoutZero := newMultiHandler()
	if fanoutZero.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected fanout with zero children to return false for Enabled")
	}
}

func TestSingleChildCollapse(t *testing.T) {
	// Single-child collapse: newMultiHandler(h) returns h itself
	h := &recordingHandler{enabledResult: true}
	result := newMultiHandler(h)

	if result != h {
		t.Error("expected newMultiHandler with single child to return the child itself")
	}
}
