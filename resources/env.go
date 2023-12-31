package resources

import (
	"github.com/caarlos0/env/v7"
	"github.com/pkg/errors"
)

type Env struct {
	DbDSN            string `env:"DB_DSN,required"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

func (r *Resource) initEnv() error {
	r.Env = Env{}
	err := env.Parse(&r.Env)
	if err != nil {
		return errors.Wrap(err, "parse env")
	}
	return nil
}
