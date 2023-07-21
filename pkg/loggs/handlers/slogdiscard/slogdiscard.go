package slogdiscard

import (
	"context"

	"golang.org/x/exp/slog"
)

/*
slog discard is required only for testing.
the main idea is to avoid writing messages to log during tests
as we don't want the messages to be saved to log while testing,
because they will not have much value
*/
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Simply ignoring log record
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Returning the same handler as far as there are no attributes to be saved
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Returning the same handler as far as there is no a group to be saved
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Always returns false as far as log record is ignored
	return false
}
