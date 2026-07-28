package project

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrProjectAccessDenied       = errors.New("project access denied")
	ErrInvalidProjectMember      = errors.New("invalid project member")
	ErrProjectSlugAlreadyExists  = errors.New("project slug already exists")
	ErrProjectKeyAlreadyExists   = errors.New("project key already exists")
	ErrInvalidProjectKey         = errors.New("invalid project key")
	ErrProjectKeyCannotBeChanged = errors.New("project key cannot be changed")
	ErrInvalidProjectDateRange   = errors.New("invalid project date range")
	ErrInvalidRepositoryURL      = errors.New("invalid repository url")
)

var projectKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]{1,9}$`)

type repositoryInterface interface {
	HasProjectAccess(
		ctx context.Context,
		projectID uint64,
		userID uint64,
		companyID uint64,
	) (bool, error)

	AreActiveCompanyMembers(
		ctx context.Context,
		companyID uint64,
		memberIDs []uint64,
	) (bool, error)

	SlugExists(
		ctx context.Context,
		companyID uint64,
		slug string,
		excludeProjectID uint64,
	) (bool, error)

	ProjectKeyExists(
		ctx context.Context,
		companyID uint64,
		projectKey string,
		excludeProjectID uint64,
	) (bool, error)

	Create(
		ctx context.Context,
		params CreateParams,
	) (*CreateResult, error)

	FindAllByUser(
		ctx context.Context,
		companyID uint64,
		userID uint64,
	) ([]Response, error)

	FindByID(
		ctx context.Context,
		projectID uint64,
		companyID uint64,
	) (*Response, error)

	Update(
		ctx context.Context,
		params UpdateParams,
	) error

	Delete(
		ctx context.Context,
		projectID uint64,
		companyID uint64,
	) error
}

type Service struct {
	repository repositoryInterface
}

func NewService(repository repositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	companyID uint64,
	userID uint64,
	request CreateRequest,
) (*Response, error) {
	memberIDs := normalizeMemberIDs(request.MemberIDs, userID)

	validMembers, err := s.repository.AreActiveCompanyMembers(
		ctx,
		companyID,
		memberIDs,
	)
	if err != nil {
		return nil, err
	}

	if !validMembers {
		return nil, ErrInvalidProjectMember
	}

	slug := strings.TrimSpace(request.Slug)

	slugExists, err := s.repository.SlugExists(
		ctx,
		companyID,
		slug,
		0,
	)
	if err != nil {
		return nil, err
	}

	if slugExists {
		return nil, ErrProjectSlugAlreadyExists
	}

	projectKey := generateProjectKeyFromSlug(slug)

	if !isValidProjectKey(projectKey) {
		return nil, ErrInvalidProjectKey
	}

	projectKeyExists, err := s.repository.ProjectKeyExists(
		ctx,
		companyID,
		projectKey,
		0,
	)
	if err != nil {
		return nil, err
	}

	if projectKeyExists {
		return nil, ErrProjectKeyAlreadyExists
	}

	repositoryURL := normalizeRepositoryURL(request.RepositoryURL)

	if repositoryURL == "" {
		return nil, ErrInvalidRepositoryURL
	}

	if !isValidDateRange(request.StartDate, request.EndDate) {
		return nil, ErrInvalidProjectDateRange
	}

	result, err := s.repository.Create(
		ctx,
		CreateParams{
			CompanyID:     companyID,
			CreatedByID:   userID,
			Name:          strings.TrimSpace(request.Name),
			Slug:          slug,
			ProjectKey:    projectKey,
			RepositoryURL: repositoryURL,
			Description:   normalizeOptionalString(request.Description),
			Status:        request.Status,
			StartDate:     normalizeOptionalString(request.StartDate),
			EndDate:       normalizeOptionalString(request.EndDate),
			MemberIDs:     memberIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	return s.repository.FindByID(
		ctx,
		result.ProjectID,
		companyID,
	)
}

func (s *Service) List(
	ctx context.Context,
	companyID uint64,
	userID uint64,
) ([]Response, error) {
	return s.repository.FindAllByUser(
		ctx,
		companyID,
		userID,
	)
}

func (s *Service) Get(
	ctx context.Context,
	projectID uint64,
	companyID uint64,
	userID uint64,
) (*Response, error) {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		userID,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, ErrProjectAccessDenied
	}

	return s.repository.FindByID(
		ctx,
		projectID,
		companyID,
	)
}

func (s *Service) Update(
	ctx context.Context,
	projectID uint64,
	companyID uint64,
	userID uint64,
	request UpdateRequest,
) (*Response, error) {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		userID,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, ErrProjectAccessDenied
	}

	currentProject, err := s.repository.FindByID(
		ctx,
		projectID,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	memberIDs := normalizeMemberIDs(request.MemberIDs, userID)

	validMembers, err := s.repository.AreActiveCompanyMembers(
		ctx,
		companyID,
		memberIDs,
	)
	if err != nil {
		return nil, err
	}

	if !validMembers {
		return nil, ErrInvalidProjectMember
	}

	slug := strings.TrimSpace(request.Slug)

	slugExists, err := s.repository.SlugExists(
		ctx,
		companyID,
		slug,
		projectID,
	)
	if err != nil {
		return nil, err
	}

	if slugExists {
		return nil, ErrProjectSlugAlreadyExists
	}

	projectKey, err := s.resolveProjectKeyForUpdate(
		ctx,
		companyID,
		projectID,
		currentProject.ProjectKey,
		request.ProjectKey,
	)
	if err != nil {
		return nil, err
	}

	repositoryURL := normalizeRepositoryURL(request.RepositoryURL)

	if repositoryURL == "" {
		return nil, ErrInvalidRepositoryURL
	}

	if !isValidDateRange(request.StartDate, request.EndDate) {
		return nil, ErrInvalidProjectDateRange
	}

	if err := s.repository.Update(
		ctx,
		UpdateParams{
			ProjectID:     projectID,
			CompanyID:     companyID,
			UpdatedByID:   userID,
			Name:          strings.TrimSpace(request.Name),
			Slug:          slug,
			ProjectKey:    projectKey,
			RepositoryURL: repositoryURL,
			Description:   normalizeOptionalString(request.Description),
			Status:        request.Status,
			StartDate:     normalizeOptionalString(request.StartDate),
			EndDate:       normalizeOptionalString(request.EndDate),
			MemberIDs:     memberIDs,
		},
	); err != nil {
		return nil, err
	}

	return s.repository.FindByID(
		ctx,
		projectID,
		companyID,
	)
}

func (s *Service) Delete(
	ctx context.Context,
	projectID uint64,
	companyID uint64,
	userID uint64,
) error {
	hasAccess, err := s.repository.HasProjectAccess(
		ctx,
		projectID,
		userID,
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
		projectID,
		companyID,
	); err != nil {
		return err
	}

	return s.repository.Delete(
		ctx,
		projectID,
		companyID,
	)
}

func (s *Service) resolveProjectKeyForUpdate(
	ctx context.Context,
	companyID uint64,
	projectID uint64,
	currentProjectKey *string,
	requestProjectKey *string,
) (*string, error) {
	if requestProjectKey == nil {
		return nil, nil
	}

	normalizedProjectKey := normalizeProjectKey(*requestProjectKey)

	if !isValidProjectKey(normalizedProjectKey) {
		return nil, ErrInvalidProjectKey
	}

	if currentProjectKey != nil {
		if *currentProjectKey != normalizedProjectKey {
			return nil, ErrProjectKeyCannotBeChanged
		}

		return nil, nil
	}

	projectKeyExists, err := s.repository.ProjectKeyExists(
		ctx,
		companyID,
		normalizedProjectKey,
		projectID,
	)
	if err != nil {
		return nil, err
	}

	if projectKeyExists {
		return nil, ErrProjectKeyAlreadyExists
	}

	return &normalizedProjectKey, nil
}

func generateProjectKeyFromSlug(slug string) string {
	normalizedSlug := strings.ToUpper(
		strings.TrimSpace(slug),
	)

	var builder strings.Builder

	for _, character := range normalizedSlug {
		isUppercaseLetter :=
			character >= 'A' &&
				character <= 'Z'

		isNumber :=
			character >= '0' &&
				character <= '9'

		if !isUppercaseLetter && !isNumber {
			continue
		}

		builder.WriteRune(character)

		if builder.Len() >= 10 {
			break
		}
	}

	return builder.String()
}

func normalizeProjectKey(projectKey string) string {
	return strings.ToUpper(
		strings.TrimSpace(projectKey),
	)
}

func normalizeRepositoryURL(repositoryURL string) string {
	return strings.TrimSpace(repositoryURL)
}

func isValidProjectKey(projectKey string) bool {
	return projectKeyPattern.MatchString(projectKey)
}

func normalizeMemberIDs(
	memberIDs []uint64,
	currentUserID uint64,
) []uint64 {
	normalizedMemberIDs := make([]uint64, 0, len(memberIDs)+1)
	seenMemberIDs := make(map[uint64]struct{}, len(memberIDs)+1)

	normalizedMemberIDs = append(normalizedMemberIDs, currentUserID)
	seenMemberIDs[currentUserID] = struct{}{}

	for _, memberID := range memberIDs {
		if memberID == 0 {
			continue
		}

		if _, exists := seenMemberIDs[memberID]; exists {
			continue
		}

		normalizedMemberIDs = append(normalizedMemberIDs, memberID)
		seenMemberIDs[memberID] = struct{}{}
	}

	return normalizedMemberIDs
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	normalizedValue := strings.TrimSpace(*value)
	if normalizedValue == "" {
		return nil
	}

	return &normalizedValue
}

func isValidDateRange(
	startDate *string,
	endDate *string,
) bool {
	normalizedStartDate := normalizeOptionalString(startDate)
	normalizedEndDate := normalizeOptionalString(endDate)

	if normalizedStartDate == nil || normalizedEndDate == nil {
		return true
	}

	start, err := time.Parse("2006-01-02", *normalizedStartDate)
	if err != nil {
		return false
	}

	end, err := time.Parse("2006-01-02", *normalizedEndDate)
	if err != nil {
		return false
	}

	return !end.Before(start)
}
