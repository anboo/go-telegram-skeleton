package resources

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Resource struct {
	Env Env
	DB  *gorm.DB
}

func Init(ctx context.Context) (*Resource, error) {
	r := &Resource{}

	err := r.initEnv()
	if err != nil {
		return nil, errors.Wrap(err, "init env")
	}

	err = r.initDb()
	if err != nil {
		return nil, errors.Wrap(err, "init db")
	}

	return r, nil
}
