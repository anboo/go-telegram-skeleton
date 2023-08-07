package subscribe

import (
	"context"
	"time"

	"telegram_bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Handler struct {
	bot *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI) *Handler {
	return &Handler{bot: bot}
}

func (h Handler) String() string {
	return "subscribe"
}

func (h Handler) Support(r *handler.BotRequest) bool {
	return !r.IsNewUser
}

func (h Handler) Handle(ctx context.Context, r *handler.BotRequest) error {
	time.Sleep(3 * time.Second)
	msg := tgbotapi.NewMessage(r.FromID, "Отправьте ссылку на RSS ленту")

	_, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "fail send welcome message")
	}

	return nil
}
