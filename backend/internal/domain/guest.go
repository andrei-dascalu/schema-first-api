package domain

import (
	"time"

	"github.com/google/uuid"
)

type RsvpStatus string

const (
	RsvpPending  RsvpStatus = "pending"
	RsvpAccepted RsvpStatus = "accepted"
	RsvpDeclined RsvpStatus = "declined"
)

type Guest struct {
	ID         uuid.UUID
	EventID    uuid.UUID
	Name       string
	Email      *string
	RsvpStatus RsvpStatus
	CreatedAt  time.Time
}
