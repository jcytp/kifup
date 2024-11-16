package log

import (
	"context"
	"log/slog"
	"os"
)

type key string

const (
	AccountIDLogKey key = "account_id"
	KifuIDLogKey    key = "kifu_id"
)

type LogHandler struct {
	slog.Handler
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if accountID, ok := ctx.Value(AccountIDLogKey).(string); ok {
		r.AddAttrs(slog.String("aid", accountID))
	}
	if kifuID, ok := ctx.Value(KifuIDLogKey).(string); ok {
		r.AddAttrs(slog.String("kid", kifuID))
	}
	return h.Handler.Handle(ctx, r)
}

func NewLogHandler(level slog.Level) *LogHandler {
	return &LogHandler{
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}),
	}
}

func Setup(level slog.Level) {
	logger := slog.New(NewLogHandler(level))
	slog.SetDefault(logger)
}
