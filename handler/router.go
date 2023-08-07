package handler

import (
	"context"
	"time"

	"telegram_bot/metrics"
	"telegram_bot/models/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type BotRequest struct {
	FromID int64

	Update    tgbotapi.Update
	User      user.User
	IsNewUser bool

	StopPropagation bool
}

type Router struct {
	handlers []Handler
}

func NewRouter(handlers ...Handler) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) Handle(ctx context.Context, request BotRequest) error {
	var found bool

	for _, h := range r.handlers {
		if !h.Support(&request) {
			continue
		}

		var err error

		startedAt := time.Now()
		found = true

		err = h.Handle(ctx, &request)
		if err != nil {
			metrics.TelegramHandlerDuration.WithLabelValues(h.String(), "error").Observe(time.Since(startedAt).Seconds())
			slog.Error("handler error", slog.Int64("from_id", request.FromID), slog.String("handler", h.String()), slog.Any("update", request))
			return nil
		} else {
			metrics.TelegramHandlerDuration.WithLabelValues(h.String(), "success").Observe(time.Since(startedAt).Seconds())
		}

		if request.StopPropagation {
			return nil
		}
	}

	if !found {
		return errors.New("not found handler")
	}
	return nil
}
