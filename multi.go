package glog

import (
	"context"
	"errors"
	"log/slog"
)

// multiHandler is a fanout handler that dispatches records to multiple child handlers.
// It ensures records are cloned per child to prevent mutation or retention issues.
type multiHandler struct {
	children []slog.Handler
}

// newMultiHandler creates a fanout handler that dispatches to multiple child handlers.
// It skips any nil handlers. If exactly one child remains after filtering, it returns
// that child directly. If zero children remain, it returns a multiHandler with no children.
// Otherwise it returns a new multiHandler wrapping all children.
func newMultiHandler(hs ...slog.Handler) slog.Handler {
	var children []slog.Handler
	for _, h := range hs {
		if h != nil {
			children = append(children, h)
		}
	}

	// If only one child, return it directly
	if len(children) == 1 {
		return children[0]
	}

	// Return fanout handler (even if zero children)
	return &multiHandler{
		children: children,
	}
}

// Enabled reports whether any child handler is enabled for the given level.
// Returns false if there are no children.
func (mh *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, child := range mh.children {
		if child.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle dispatches the record to all children that are enabled for its level.
// Each child receives a cloned record to prevent mutation or retention issues.
// All children are called even if one returns an error; errors are joined and returned.
func (mh *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, child := range mh.children {
		if !child.Enabled(ctx, r.Level) {
			continue
		}
		err := child.Handle(ctx, r.Clone())
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// WithAttrs returns a new multiHandler whose children have received the given attributes.
// If attrs is empty, returns the receiver unchanged.
// Creates a brand new slice of handlers to avoid aliasing the original.
func (mh *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return mh
	}

	newChildren := make([]slog.Handler, len(mh.children))
	for i, child := range mh.children {
		newChildren[i] = child.WithAttrs(attrs)
	}

	return &multiHandler{
		children: newChildren,
	}
}

// WithGroup returns a new multiHandler whose children have received the given group.
// If name is empty, returns the receiver unchanged.
// Creates a brand new slice of handlers to avoid aliasing the original.
func (mh *multiHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return mh
	}

	newChildren := make([]slog.Handler, len(mh.children))
	for i, child := range mh.children {
		newChildren[i] = child.WithGroup(name)
	}

	return &multiHandler{
		children: newChildren,
	}
}
