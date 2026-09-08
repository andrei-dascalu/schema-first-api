package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID
	Name      string
	Date      time.Time
	Location  *string
	CreatedAt time.Time
}
