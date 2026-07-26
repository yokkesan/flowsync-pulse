package realtime

import (
	"context"

	"flowsync-pulse/backend/internal/extension"
)

type ExtensionNotifier struct {
	hub *Hub
}

func NewExtensionNotifier(
	hub *Hub,
) *ExtensionNotifier {
	return &ExtensionNotifier{
		hub: hub,
	}
}

func (n *ExtensionNotifier) NotifyWorkContextChanged(
	ctx context.Context,
	companyID uint64,
	notification extension.WorkContextNotification,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if n.hub == nil {
		return nil
	}

	return n.hub.BroadcastWorkContextChanged(
		companyID,
		WorkContextChangedPayload{
			UserID:          notification.UserID,
			ProjectID:       notification.ProjectID,
			ProjectName:     notification.ProjectName,
			TaskID:          notification.TaskID,
			TaskKey:         notification.TaskKey,
			TaskName:        notification.TaskName,
			RepositoryName:  notification.RepositoryName,
			BranchName:      notification.BranchName,
			TicketKey:       notification.TicketKey,
			WorkspaceName:   notification.WorkspaceName,
			MatchStatus:     notification.MatchStatus,
			SessionStatus:   notification.SessionStatus,
			ExtensionActive: notification.ExtensionActive,
			StartedAt:       notification.StartedAt,
			LastHeartbeatAt: notification.LastHeartbeatAt,
			EndedAt:         notification.EndedAt,
			EndReason:       notification.EndReason,
		},
	)
}
