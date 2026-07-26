package realtime

import (
	"context"
	"errors"
	"time"
)

var (
	ErrConnectionTicketNotFound = errors.New(
		"connection ticket not found",
	)
	ErrConnectionTicketExpired = errors.New(
		"connection ticket expired",
	)
	ErrConnectionTicketUsed = errors.New(
		"connection ticket already used",
	)
)

type ConnectionTicket struct {
	Token     string
	UserID    uint64
	CompanyID uint64
	ExpiresAt time.Time
	UsedAt    *time.Time
}

type ConnectionTicketStore interface {
	Create(
		ctx context.Context,
		ticket ConnectionTicket,
	) error

	Consume(
		ctx context.Context,
		token string,
		now time.Time,
	) (ConnectionTicket, error)

	DeleteExpired(
		ctx context.Context,
		now time.Time,
	) error
}
