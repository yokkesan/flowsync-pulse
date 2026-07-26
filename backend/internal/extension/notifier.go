package extension

import (
	"context"
	"time"
)

type WorkContextNotification struct {
	UserID          uint64
	ProjectID       *uint64
	ProjectName     *string
	TaskID          *uint64
	TaskKey         *string
	TaskName        *string
	RepositoryName  *string
	BranchName      *string
	TicketKey       *string
	WorkspaceName   *string
	MatchStatus     *string
	SessionStatus   *string
	ExtensionActive bool
	StartedAt       *time.Time
	LastHeartbeatAt *time.Time
	EndedAt         *time.Time
	EndReason       *string
}

type WorkContextNotifier interface {
	NotifyWorkContextChanged(
		ctx context.Context,
		companyID uint64,
		notification WorkContextNotification,
	) error
}

type NoopWorkContextNotifier struct{}

func NewNoopWorkContextNotifier() *NoopWorkContextNotifier {
	return &NoopWorkContextNotifier{}
}

func (n *NoopWorkContextNotifier) NotifyWorkContextChanged(
	ctx context.Context,
	companyID uint64,
	notification WorkContextNotification,
) error {
	return nil
}
