package realtime

import (
	"context"
	"sync"
	"time"
)

type MemoryConnectionTicketStore struct {
	mu sync.Mutex

	tickets map[string]ConnectionTicket
}

func NewMemoryConnectionTicketStore() *MemoryConnectionTicketStore {
	return &MemoryConnectionTicketStore{
		tickets: make(
			map[string]ConnectionTicket,
		),
	}
}

func (s *MemoryConnectionTicketStore) Create(
	ctx context.Context,
	ticket ConnectionTicket,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tickets[ticket.Token] = ticket

	return nil
}

func (s *MemoryConnectionTicketStore) Consume(
	ctx context.Context,
	token string,
	now time.Time,
) (ConnectionTicket, error) {
	if err := ctx.Err(); err != nil {
		return ConnectionTicket{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.tickets[token]
	if !exists {
		return ConnectionTicket{},
			ErrConnectionTicketNotFound
	}

	if ticket.UsedAt != nil {
		delete(
			s.tickets,
			token,
		)

		return ConnectionTicket{},
			ErrConnectionTicketUsed
	}

	if !ticket.ExpiresAt.After(now) {
		delete(
			s.tickets,
			token,
		)

		return ConnectionTicket{},
			ErrConnectionTicketExpired
	}

	usedAt := now
	ticket.UsedAt = &usedAt

	// 1回限りのチケットなので、利用時点で削除する。
	delete(
		s.tickets,
		token,
	)

	return ticket, nil
}

func (s *MemoryConnectionTicketStore) DeleteExpired(
	ctx context.Context,
	now time.Time,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for token, ticket := range s.tickets {
		if ticket.UsedAt != nil ||
			!ticket.ExpiresAt.After(now) {
			delete(
				s.tickets,
				token,
			)
		}
	}

	return nil
}
