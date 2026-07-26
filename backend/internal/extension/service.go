package extension

import (
	"context"
	"errors"
	"log"
	"net/url"
	"path"
	"strings"
	"time"
)

var (
	ErrInvalidRepositoryURL = errors.New(
		"invalid repository url",
	)
	ErrExtensionProjectAccessDenied = errors.New(
		"extension project access denied",
	)
	ErrInvalidOccurredAt = errors.New(
		"invalid occurred at",
	)
)

type repositoryInterface interface {
	FindRepositoryCandidates(
		ctx context.Context,
		companyID uint64,
	) ([]RepositoryCandidate, error)

	HasProjectAccess(
		ctx context.Context,
		projectID uint64,
		userID uint64,
		companyID uint64,
	) (bool, error)

	FindTaskByKey(
		ctx context.Context,
		projectID uint64,
		taskKey string,
	) (*TaskMatch, error)

	FindTaskByBranch(
		ctx context.Context,
		projectID uint64,
		branchName string,
	) (*TaskMatch, error)

	SaveWorkContext(
		ctx context.Context,
		params StartSessionParams,
	) (SessionUpdateResult, error)

	Heartbeat(
		ctx context.Context,
		userID uint64,
		occurredAt time.Time,
	) (HeartbeatResult, error)

	Disconnect(
		ctx context.Context,
		userID uint64,
		occurredAt time.Time,
	) (*DisconnectResult, error)
}

type Service struct {
	repository repositoryInterface
	notifier   WorkContextNotifier
}

func NewService(
	repository repositoryInterface,
	notifier WorkContextNotifier,
) *Service {
	if notifier == nil {
		notifier = NewNoopWorkContextNotifier()
	}

	return &Service{
		repository: repository,
		notifier:   notifier,
	}
}

func (s *Service) WorkContext(
	ctx context.Context,
	userID uint64,
	companyID uint64,
	request WorkContextRequest,
) (WorkContextResponse, error) {
	repositoryName := strings.TrimSpace(
		request.RepositoryName,
	)
	repositoryURL := strings.TrimSpace(
		request.RepositoryURL,
	)
	branchName := strings.TrimSpace(
		request.BranchName,
	)
	ticketKey := normalizeOptionalString(
		request.TicketKey,
	)
	workspaceName := normalizeOptionalString(
		request.WorkspaceName,
	)
	extensionVersion := normalizeOptionalString(
		request.ExtensionVersion,
	)

	occurredAt, err := normalizeOccurredAt(
		request.OccurredAt,
	)
	if err != nil {
		return WorkContextResponse{}, err
	}

	normalizedRepositoryURL, err :=
		normalizeRepositoryURL(repositoryURL)
	if err != nil {
		return WorkContextResponse{},
			ErrInvalidRepositoryURL
	}

	repository, err := s.findRepository(
		ctx,
		companyID,
		normalizedRepositoryURL,
	)
	if err != nil {
		return WorkContextResponse{}, err
	}

	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		repository.ProjectID,
		userID,
		companyID,
	)
	if err != nil {
		return WorkContextResponse{}, err
	}

	if !hasAccess {
		return WorkContextResponse{},
			ErrExtensionProjectAccessDenied
	}

	taskMatch, matchStatus, err :=
		s.resolveTaskMatch(
			ctx,
			repository.ProjectID,
			ticketKey,
			branchName,
		)
	if err != nil {
		return WorkContextResponse{}, err
	}

	var taskID *uint64
	var taskKey *string
	var taskName *string

	if taskMatch != nil {
		taskIDValue := taskMatch.TaskID
		taskKeyValue := taskMatch.TaskKey
		taskNameValue := taskMatch.TaskName

		taskID = &taskIDValue
		taskKey = &taskKeyValue
		taskName = &taskNameValue
	}

	session, err := s.repository.SaveWorkContext(
		ctx,
		StartSessionParams{
			UserID:        userID,
			ProjectID:     repository.ProjectID,
			RepositoryID:  repository.RepositoryID,
			TaskID:        taskID,
			BranchName:    branchName,
			TicketKey:     ticketKey,
			WorkspaceName: workspaceName,
			MatchStatus:   matchStatus,
			OccurredAt:    occurredAt,
		},
	)
	if err != nil {
		return WorkContextResponse{}, err
	}

	response := WorkContextResponse{
		SessionID:        session.SessionID,
		ProjectID:        repository.ProjectID,
		ProjectName:      repository.ProjectName,
		TaskID:           taskID,
		TaskKey:          taskKey,
		TaskName:         taskName,
		RepositoryName:   repositoryName,
		RepositoryURL:    repositoryURL,
		BranchName:       branchName,
		TicketKey:        ticketKey,
		WorkspaceName:    workspaceName,
		ExtensionVersion: extensionVersion,
		MatchStatus:      matchStatus,
		Status:           session.Status,
		StartedAt:        session.StartedAt,
		LastHeartbeatAt:  session.LastHeartbeatAt,
		EndedAt:          session.EndedAt,
	}

	projectID := repository.ProjectID
	projectName := repository.ProjectName
	sessionStatus := session.Status
	startedAt := session.StartedAt
	lastHeartbeatAt := session.LastHeartbeatAt
	repositoryNameValue := repositoryName
	branchNameValue := branchName
	matchStatusValue := matchStatus

	if err := s.notifier.NotifyWorkContextChanged(
		ctx,
		companyID,
		WorkContextNotification{
			UserID:          userID,
			ProjectID:       &projectID,
			ProjectName:     &projectName,
			TaskID:          taskID,
			TaskKey:         taskKey,
			TaskName:        taskName,
			RepositoryName:  &repositoryNameValue,
			BranchName:      &branchNameValue,
			TicketKey:       ticketKey,
			WorkspaceName:   workspaceName,
			MatchStatus:     &matchStatusValue,
			SessionStatus:   &sessionStatus,
			ExtensionActive: true,
			StartedAt:       &startedAt,
			LastHeartbeatAt: &lastHeartbeatAt,
			EndedAt:         session.EndedAt,
			EndReason:       nil,
		},
	); err != nil {
		log.Printf(
			"failed to notify realtime work context change: user_id=%d company_id=%d error=%v",
			userID,
			companyID,
			err,
		)
	}

	return response, nil
}

func (s *Service) Heartbeat(
	ctx context.Context,
	userID uint64,
	companyID uint64,
	request HeartbeatRequest,
) (HeartbeatResponse, error) {
	occurredAt, err := normalizeOccurredAt(
		request.OccurredAt,
	)
	if err != nil {
		return HeartbeatResponse{}, err
	}

	result, err := s.repository.Heartbeat(
		ctx,
		userID,
		occurredAt,
	)
	if err != nil {
		return HeartbeatResponse{}, err
	}

	response := HeartbeatResponse{
		SessionID:       result.SessionID,
		Status:          result.Status,
		LastHeartbeatAt: result.LastHeartbeatAt,
	}

	sessionStatus := result.Status
	lastHeartbeatAt := result.LastHeartbeatAt

	if err := s.notifier.NotifyWorkContextChanged(
		ctx,
		companyID,
		WorkContextNotification{
			UserID:          userID,
			SessionStatus:   &sessionStatus,
			ExtensionActive: true,
			LastHeartbeatAt: &lastHeartbeatAt,
		},
	); err != nil {
		log.Printf(
			"failed to notify realtime heartbeat change: user_id=%d company_id=%d error=%v",
			userID,
			companyID,
			err,
		)
	}

	return response, nil
}

func (s *Service) Disconnect(
	ctx context.Context,
	userID uint64,
	companyID uint64,
	request DisconnectRequest,
) (*DisconnectResponse, error) {
	occurredAt, err := normalizeOccurredAt(
		request.OccurredAt,
	)
	if err != nil {
		return nil, err
	}

	result, err := s.repository.Disconnect(
		ctx,
		userID,
		occurredAt,
	)
	if err != nil {
		return nil, err
	}

	// 既に終了済みの場合も正常終了とする。
	if result == nil {
		return nil, nil
	}

	endReason := result.EndReason
	endedAt := result.EndedAt
	sessionStatus := result.Status

	response := &DisconnectResponse{
		SessionID: result.SessionID,
		Status:    result.Status,
		EndReason: &endReason,
		EndedAt:   &endedAt,
	}

	if err := s.notifier.NotifyWorkContextChanged(
		ctx,
		companyID,
		WorkContextNotification{
			UserID:          userID,
			SessionStatus:   &sessionStatus,
			ExtensionActive: false,
			EndedAt:         &endedAt,
			EndReason:       &endReason,
		},
	); err != nil {
		log.Printf(
			"failed to notify realtime disconnect change: user_id=%d company_id=%d error=%v",
			userID,
			companyID,
			err,
		)
	}

	return response, nil
}

func (s *Service) findRepository(
	ctx context.Context,
	companyID uint64,
	normalizedRepositoryURL string,
) (RepositoryCandidate, error) {
	candidates, err :=
		s.repository.FindRepositoryCandidates(
			ctx,
			companyID,
		)
	if err != nil {
		return RepositoryCandidate{}, err
	}

	for _, candidate := range candidates {
		normalizedCandidateURL, err :=
			normalizeRepositoryURL(
				candidate.RemoteURL,
			)
		if err != nil {
			continue
		}

		if strings.EqualFold(
			normalizedCandidateURL,
			normalizedRepositoryURL,
		) {
			return candidate, nil
		}
	}

	return RepositoryCandidate{},
		ErrRepositoryNotFound
}

func (s *Service) resolveTaskMatch(
	ctx context.Context,
	projectID uint64,
	ticketKey *string,
	branchName string,
) (
	*TaskMatch,
	string,
	error,
) {
	if ticketKey == nil {
		branchTask, err :=
			s.repository.FindTaskByBranch(
				ctx,
				projectID,
				branchName,
			)
		if err != nil {
			return nil, "", err
		}

		if branchTask == nil {
			return nil,
				MatchStatusBranchNotMatched,
				nil
		}

		return nil,
			MatchStatusTicketNotFound,
			nil
	}

	taskByKey, err := s.repository.FindTaskByKey(
		ctx,
		projectID,
		*ticketKey,
	)
	if err != nil {
		return nil, "", err
	}

	if taskByKey == nil {
		return nil,
			MatchStatusTicketNotFound,
			nil
	}

	if taskByKey.BranchName != branchName {
		return nil,
			MatchStatusTicketBranchMismatch,
			nil
	}

	return taskByKey,
		MatchStatusMatched,
		nil
}

func normalizeRepositoryURL(
	rawRepositoryURL string,
) (string, error) {
	normalizedValue := strings.TrimSpace(
		rawRepositoryURL,
	)
	if normalizedValue == "" {
		return "", ErrInvalidRepositoryURL
	}

	normalizedValue = strings.TrimSuffix(
		normalizedValue,
		"/",
	)

	if isSCPStyleRepositoryURL(
		normalizedValue,
	) {
		parts := strings.SplitN(
			normalizedValue,
			":",
			2,
		)
		if len(parts) != 2 {
			return "", ErrInvalidRepositoryURL
		}

		hostPart := parts[0]
		repositoryPath := parts[1]

		if atIndex := strings.LastIndex(
			hostPart,
			"@",
		); atIndex >= 0 {
			hostPart = hostPart[atIndex+1:]
		}

		return buildNormalizedRepositoryURL(
			hostPart,
			repositoryPath,
		)
	}

	parsedURL, err := url.Parse(
		normalizedValue,
	)
	if err != nil {
		return "", ErrInvalidRepositoryURL
	}

	switch parsedURL.Scheme {
	case "http", "https", "ssh", "git":
	default:
		return "", ErrInvalidRepositoryURL
	}

	host := parsedURL.Hostname()
	if host == "" {
		return "", ErrInvalidRepositoryURL
	}

	return buildNormalizedRepositoryURL(
		host,
		parsedURL.Path,
	)
}

func isSCPStyleRepositoryURL(
	value string,
) bool {
	if strings.Contains(
		value,
		"://",
	) {
		return false
	}

	colonIndex := strings.Index(
		value,
		":",
	)

	return colonIndex > 0 &&
		colonIndex < len(value)-1
}

func buildNormalizedRepositoryURL(
	host string,
	repositoryPath string,
) (string, error) {
	normalizedHost := strings.ToLower(
		strings.TrimSpace(host),
	)
	if normalizedHost == "" {
		return "", ErrInvalidRepositoryURL
	}

	normalizedPath := strings.TrimSpace(
		repositoryPath,
	)
	normalizedPath = strings.TrimPrefix(
		normalizedPath,
		"/",
	)
	normalizedPath = strings.TrimSuffix(
		normalizedPath,
		"/",
	)
	normalizedPath = strings.TrimSuffix(
		normalizedPath,
		".git",
	)
	normalizedPath = path.Clean(
		normalizedPath,
	)

	if normalizedPath == "" ||
		normalizedPath == "." ||
		normalizedPath == "/" {
		return "", ErrInvalidRepositoryURL
	}

	return normalizedHost + "/" +
		normalizedPath, nil
}

func normalizeOptionalString(
	value *string,
) *string {
	if value == nil {
		return nil
	}

	normalizedValue := strings.TrimSpace(
		*value,
	)
	if normalizedValue == "" {
		return nil
	}

	return &normalizedValue
}

func normalizeOccurredAt(
	occurredAt time.Time,
) (time.Time, error) {
	if occurredAt.IsZero() {
		return time.Time{},
			ErrInvalidOccurredAt
	}

	return occurredAt.UTC(), nil
}
