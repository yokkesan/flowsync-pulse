package task

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrProjectAccessDenied = errors.New(
		"project access denied",
	)
	ErrAssigneeNotProjectMember = errors.New(
		"assignee is not an active project member",
	)
	ErrBranchAlreadyExists = errors.New(
		"branch already exists in project",
	)
	ErrInvalidDateRange = errors.New(
		"due date must not be earlier than start date",
	)
	ErrTaskNotFound = errors.New(
		"task not found",
	)
)

type TaskCreator interface {
	HasProjectAccess(
		ctx context.Context,
		projectID uint64,
		userID uint64,
		companyID uint64,
	) (bool, error)

	IsActiveProjectMember(
		ctx context.Context,
		projectID uint64,
		userID uint64,
	) (bool, error)

	BranchExistsInProject(
		ctx context.Context,
		projectID uint64,
		branchName string,
		excludeTaskID uint64,
	) (bool, error)

	Create(
		ctx context.Context,
		params CreateParams,
	) (CreateResult, error)

	Update(
		ctx context.Context,
		params UpdateParams,
	) error

	Delete(
		ctx context.Context,
		projectID uint64,
		taskID uint64,
	) error

	FindByID(
		ctx context.Context,
		params FindByIDParams,
	) (Response, error)

	FindAllByProjectID(
		ctx context.Context,
		projectID uint64,
	) ([]Response, error)
}

type Service struct {
	repository TaskCreator
}

func NewService(repository TaskCreator) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	projectID uint64,
	authUserID uint64,
	companyID uint64,
	request CreateRequest,
) (Response, error) {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		authUserID,
		companyID,
	)
	if err != nil {
		return Response{}, err
	}

	if !hasAccess {
		return Response{}, ErrProjectAccessDenied
	}

	isProjectMember, err := s.repository.IsActiveProjectMember(
		ctx,
		projectID,
		request.AssigneeUserID,
	)
	if err != nil {
		return Response{}, err
	}

	if !isProjectMember {
		return Response{}, ErrAssigneeNotProjectMember
	}

	name := strings.TrimSpace(request.Name)
	branchName := strings.TrimSpace(request.BranchName)

	branchExists, err := s.repository.BranchExistsInProject(
		ctx,
		projectID,
		branchName,
		0,
	)
	if err != nil {
		return Response{}, err
	}

	if branchExists {
		return Response{}, ErrBranchAlreadyExists
	}

	if err := validateDateRange(
		request.StartDate,
		request.DueDate,
	); err != nil {
		return Response{}, err
	}

	description := trimOptionalString(
		request.Description,
	)

	var completedAt *time.Time

	if request.Status == StatusCompleted {
		now := time.Now().UTC()
		completedAt = &now
	}

	result, err := s.repository.Create(
		ctx,
		CreateParams{
			ProjectID:      projectID,
			Name:           name,
			Description:    description,
			AssigneeUserID: request.AssigneeUserID,
			BranchName:     branchName,
			Status:         request.Status,
			Priority:       request.Priority,
			StartDate:      request.StartDate,
			DueDate:        request.DueDate,
			CompletedAt:    completedAt,
		},
	)
	if err != nil {
		return Response{}, err
	}

	response, err := s.repository.FindByID(
		ctx,
		FindByIDParams{
			ProjectID: projectID,
			TaskID:    result.TaskID,
		},
	)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (s *Service) List(
	ctx context.Context,
	projectID uint64,
	authUserID uint64,
	companyID uint64,
) (ListResponse, error) {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		authUserID,
		companyID,
	)
	if err != nil {
		return ListResponse{}, err
	}

	if !hasAccess {
		return ListResponse{}, ErrProjectAccessDenied
	}

	tasks, err := s.repository.FindAllByProjectID(
		ctx,
		projectID,
	)
	if err != nil {
		return ListResponse{}, err
	}

	return ListResponse{
		Tasks: tasks,
	}, nil
}

func (s *Service) Get(
	ctx context.Context,
	projectID uint64,
	taskID uint64,
	authUserID uint64,
	companyID uint64,
) (Response, error) {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		authUserID,
		companyID,
	)
	if err != nil {
		return Response{}, err
	}

	if !hasAccess {
		return Response{}, ErrProjectAccessDenied
	}

	response, err := s.repository.FindByID(
		ctx,
		FindByIDParams{
			ProjectID: projectID,
			TaskID:    taskID,
		},
	)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (s *Service) Update(
	ctx context.Context,
	projectID uint64,
	taskID uint64,
	authUserID uint64,
	companyID uint64,
	request UpdateRequest,
) (Response, error) {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		authUserID,
		companyID,
	)
	if err != nil {
		return Response{}, err
	}

	if !hasAccess {
		return Response{}, ErrProjectAccessDenied
	}

	currentTask, err := s.repository.FindByID(
		ctx,
		FindByIDParams{
			ProjectID: projectID,
			TaskID:    taskID,
		},
	)
	if err != nil {
		return Response{}, err
	}

	isProjectMember, err := s.repository.IsActiveProjectMember(
		ctx,
		projectID,
		request.AssigneeUserID,
	)
	if err != nil {
		return Response{}, err
	}

	if !isProjectMember {
		return Response{}, ErrAssigneeNotProjectMember
	}

	name := strings.TrimSpace(request.Name)
	branchName := strings.TrimSpace(request.BranchName)

	branchExists, err := s.repository.BranchExistsInProject(
		ctx,
		projectID,
		branchName,
		taskID,
	)
	if err != nil {
		return Response{}, err
	}

	if branchExists {
		return Response{}, ErrBranchAlreadyExists
	}

	if err := validateDateRange(
		request.StartDate,
		request.DueDate,
	); err != nil {
		return Response{}, err
	}

	description := trimOptionalString(
		request.Description,
	)

	completedAt := currentTask.CompletedAt

	if request.Status == StatusCompleted {
		if completedAt == nil {
			now := time.Now().UTC()
			completedAt = &now
		}
	} else {
		completedAt = nil
	}

	err = s.repository.Update(
		ctx,
		UpdateParams{
			ProjectID:      projectID,
			TaskID:         taskID,
			Name:           name,
			Description:    description,
			AssigneeUserID: request.AssigneeUserID,
			BranchName:     branchName,
			Status:         request.Status,
			Priority:       request.Priority,
			StartDate:      request.StartDate,
			DueDate:        request.DueDate,
			CompletedAt:    completedAt,
		},
	)
	if err != nil {
		return Response{}, err
	}

	response, err := s.repository.FindByID(
		ctx,
		FindByIDParams{
			ProjectID: projectID,
			TaskID:    taskID,
		},
	)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (s *Service) Delete(
	ctx context.Context,
	projectID uint64,
	taskID uint64,
	authUserID uint64,
	companyID uint64,
) error {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		authUserID,
		companyID,
	)
	if err != nil {
		return err
	}

	if !hasAccess {
		return ErrProjectAccessDenied
	}

	if _, err := s.repository.FindByID(
		ctx,
		FindByIDParams{
			ProjectID: projectID,
			TaskID:    taskID,
		},
	); err != nil {
		return err
	}

	if err := s.repository.Delete(
		ctx,
		projectID,
		taskID,
	); err != nil {
		return err
	}

	return nil
}

func validateDateRange(
	startDate *string,
	dueDate *string,
) error {
	if startDate == nil || dueDate == nil {
		return nil
	}

	start, err := time.Parse(
		time.DateOnly,
		*startDate,
	)
	if err != nil {
		return err
	}

	due, err := time.Parse(
		time.DateOnly,
		*dueDate,
	)
	if err != nil {
		return err
	}

	if due.Before(start) {
		return ErrInvalidDateRange
	}

	return nil
}

func trimOptionalString(
	value *string,
) *string {
	if value == nil {
		return nil
	}

	trimmedValue := strings.TrimSpace(*value)

	if trimmedValue == "" {
		return nil
	}

	return &trimmedValue
}
