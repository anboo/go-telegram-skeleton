package models

import (
	"time"
)

type Lot struct {
	Id          string
	ExternalId  string
	Title       string
	Description string
	CreatedAt   time.Time
}
