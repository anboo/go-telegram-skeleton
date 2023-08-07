package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"telegram_bot/db/user"
	"telegram_bot/handler"
	"telegram_bot/handler/subscribe"
	user_model "telegram_bot/models/user"
	"telegram_bot/resources"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	res, err := resources.Init(ctx)
	if err != nil {
		panic(err)
	}

	userRepository := user.NewRepository(res.DB)

	bot, err := tgbotapi.NewBotAPI(res.Env.TelegramBotToken)
	if err != nil {
		slog.Error("fail initialize bot", err)
		os.Exit(1)
	}
	slog.Info("authorized an account", slog.String("self_username", bot.Self.UserName))

	router := handler.NewRouter(
		subscribe.NewHandler(bot),
	)

	grp, _ := errgroup.WithContext(ctx)
	grp.SetLimit(100)

	ch := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Limit: 100})
	go func() {
		select {
		case <-ctx.Done():
			bot.StopReceivingUpdates()
		case update := <-ch:
			grp.Go(func() error {
				registeredUser, isNew, err := userRepository.GetOrCreate(ctx, user_model.NewUser(
					fmt.Sprintf("%d", int(update.FromChat().ID)),
					update.FromChat().UserName,
					update.FromChat().FirstName+" "+update.FromChat().LastName,
				))

				if err != nil {
					slog.Error("try get or create user", err)
					return errors.Wrap(err, "try get or create user")
				}

				err = router.Handle(context.Background(), handler.BotRequest{
					FromID:    update.FromChat().ID,
					Update:    update,
					User:      registeredUser,
					IsNewUser: isNew,
				})

				if err != nil {
					slog.Error("fail process handler",
						err,
						slog.Int64("chat_id", update.FromChat().ID),
						slog.Any("update", update))
				}

				return err
			})
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down telegram bot, waiting handlers")
	err = grp.Wait()
	slog.Info("done wait handlers", slog.Any("last_error", err))
}
