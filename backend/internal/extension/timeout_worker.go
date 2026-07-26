package extension

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultSessionTimeoutSeconds       = 90
	defaultTimeoutCheckIntervalSeconds = 30
)

type TimeoutWorker struct {
	db             *sql.DB
	notifier       WorkContextNotifier
	sessionTimeout time.Duration
	checkInterval  time.Duration
}

type timedOutSession struct {
	SessionID uint64
	UserID    uint64
	CompanyID uint64
}

func NewTimeoutWorker(
	db *sql.DB,
	notifier WorkContextNotifier,
) *TimeoutWorker {
	if notifier == nil {
		notifier = NewNoopWorkContextNotifier()
	}

	return &TimeoutWorker{
		db:       db,
		notifier: notifier,
		sessionTimeout: time.Duration(
			getPositiveEnvironmentInteger(
				"EXTENSION_SESSION_TIMEOUT_SECONDS",
				defaultSessionTimeoutSeconds,
			),
		) * time.Second,
		checkInterval: time.Duration(
			getPositiveEnvironmentInteger(
				"EXTENSION_TIMEOUT_CHECK_INTERVAL_SECONDS",
				defaultTimeoutCheckIntervalSeconds,
			),
		) * time.Second,
	}
}

func (w *TimeoutWorker) Start(
	ctx context.Context,
) {
	if w.db == nil {
		log.Print(
			"extension timeout worker was not started: database is nil",
		)
		return
	}

	log.Printf(
		"extension timeout worker started: session_timeout=%s check_interval=%s",
		w.sessionTimeout,
		w.checkInterval,
	)

	ticker := time.NewTicker(
		w.checkInterval,
	)
	defer ticker.Stop()

	w.processTimeouts(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Print(
				"extension timeout worker stopped",
			)
			return

		case <-ticker.C:
			w.processTimeouts(ctx)
		}
	}
}

func (w *TimeoutWorker) processTimeouts(
	ctx context.Context,
) {
	if err := ctx.Err(); err != nil {
		return
	}

	checkedAt := time.Now().UTC()
	timeoutThreshold := checkedAt.Add(
		-w.sessionTimeout,
	)

	sessions, err := w.findTimeoutCandidates(
		ctx,
		timeoutThreshold,
	)
	if err != nil {
		log.Printf(
			"failed to find extension timeout candidates: error=%v",
			err,
		)
		return
	}

	timedOutCount := 0

	for _, session := range sessions {
		updated, err := w.markSessionTimedOut(
			ctx,
			session.SessionID,
			timeoutThreshold,
			checkedAt,
		)
		if err != nil {
			log.Printf(
				"failed to time out extension session: session_id=%d user_id=%d company_id=%d error=%v",
				session.SessionID,
				session.UserID,
				session.CompanyID,
				err,
			)
			continue
		}

		// 候補取得後にheartbeatなどで更新されていた場合は、
		// activeのままとし、通知も送信しない。
		if !updated {
			continue
		}

		timedOutCount++

		sessionStatus := SessionStatusTimedOut
		endReason := EndReasonTimeout
		endedAt := checkedAt

		if err := w.notifier.NotifyWorkContextChanged(
			ctx,
			session.CompanyID,
			WorkContextNotification{
				UserID:          session.UserID,
				SessionStatus:   &sessionStatus,
				ExtensionActive: false,
				EndedAt:         &endedAt,
				EndReason:       &endReason,
			},
		); err != nil {
			// DB上のタイムアウト処理は完了しているため、
			// Realtime通知失敗では処理を失敗扱いにしない。
			log.Printf(
				"failed to notify realtime session timeout: session_id=%d user_id=%d company_id=%d error=%v",
				session.SessionID,
				session.UserID,
				session.CompanyID,
				err,
			)
		}
	}

	if timedOutCount > 0 {
		log.Printf(
			"extension sessions timed out: count=%d checked_at=%s",
			timedOutCount,
			checkedAt.Format(time.RFC3339),
		)
	}
}

func (w *TimeoutWorker) findTimeoutCandidates(
	ctx context.Context,
	timeoutThreshold time.Time,
) ([]timedOutSession, error) {
	rows, err := w.db.QueryContext(
		ctx,
		`
			SELECT
				ws.id,
				ws.user_id,
				u.company_id
			FROM work_sessions AS ws
			INNER JOIN users AS u
				ON u.id = ws.user_id
			WHERE ws.status = ?
			  AND ws.last_heartbeat_at < ?
			ORDER BY ws.id ASC
		`,
		SessionStatusActive,
		timeoutThreshold,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make(
		[]timedOutSession,
		0,
	)

	for rows.Next() {
		var session timedOutSession

		if err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&session.CompanyID,
		); err != nil {
			return nil, err
		}

		sessions = append(
			sessions,
			session,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (w *TimeoutWorker) markSessionTimedOut(
	ctx context.Context,
	sessionID uint64,
	timeoutThreshold time.Time,
	endedAt time.Time,
) (bool, error) {
	result, err := w.db.ExecContext(
		ctx,
		`
			UPDATE work_sessions
			SET
				status = ?,
				ended_at = ?,
				end_reason = ?,
				updated_at = ?
			WHERE id = ?
			  AND status = ?
			  AND last_heartbeat_at < ?
		`,
		SessionStatusTimedOut,
		endedAt,
		EndReasonTimeout,
		endedAt,
		sessionID,
		SessionStatusActive,
		timeoutThreshold,
	)
	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affectedRows > 0, nil
}

func getPositiveEnvironmentInteger(
	name string,
	defaultValue int,
) int {
	rawValue := os.Getenv(name)
	if rawValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(
		rawValue,
	)
	if err != nil || value <= 0 {
		log.Printf(
			"invalid environment value, using default: name=%s value=%q default=%d",
			name,
			rawValue,
			defaultValue,
		)
		return defaultValue
	}

	return value
}
