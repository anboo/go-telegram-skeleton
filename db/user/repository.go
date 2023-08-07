package user

import (
	"context"
	"time"

	"telegram_bot/models/user"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetOrCreate(ctx context.Context, newUser user.User) (registered user.User, isNewUser bool, err error) {
	newUser.ID = uuid.New().String()
	newUser.CreatedAt = time.Now()

	err = r.db.WithContext(ctx).Create(newUser).Error
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		err = r.db.WithContext(ctx).First(&registered, "telegram_id = ?", newUser.TelegramID).Error
		if err != nil {
			return user.User{}, false, errors.Wrap(err, "try fetch actual user after conflict update")
		}
		return registered, false, nil
	case err != nil:
		return user.User{}, false, errors.Wrap(err, "try create user")
	}

	return newUser, true, nil
}
