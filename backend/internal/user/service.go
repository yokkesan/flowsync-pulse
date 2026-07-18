package user

import ( "context" "errors" "regexp" "strings" "golang.org/x/crypto/bcrypt" )

var (
	ErrPasswordMismatch   = errors.New("password confirmation does not match")
	ErrInvalidCompanySlug = errors.New("company slug must contain only lowercase letters, numbers, and hyphens")
)

var companySlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type CompanyOwnerCreator interface {
	CreateCompanyOwner(
		ctx context.Context,
		params CreateCompanyOwnerParams,
	) (CreateCompanyOwnerResult, error)
}

type Service struct {
	repository CompanyOwnerCreator
}

func NewService(repository CompanyOwnerCreator) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) RegisterCompanyOwner(
	ctx context.Context,
	request RegisterRequest,
) (RegisterResponse, error) {
	companyName := strings.TrimSpace(request.CompanyName)
	companySlug := strings.ToLower(strings.TrimSpace(request.CompanySlug))
	displayName := strings.TrimSpace(request.DisplayName)
	email := strings.ToLower(strings.TrimSpace(request.Email))
	avatarKey := strings.TrimSpace(request.AvatarKey)

	if request.Password != request.PasswordConfirm {
		return RegisterResponse{}, ErrPasswordMismatch
	}

	if !companySlugPattern.MatchString(companySlug) {
		return RegisterResponse{}, ErrInvalidCompanySlug
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return RegisterResponse{}, err
	}

	result, err := s.repository.CreateCompanyOwner(
		ctx,
		CreateCompanyOwnerParams{
			CompanyName:  companyName,
			CompanySlug:  companySlug,
			DisplayName:  displayName,
			Email:        email,
			PasswordHash: string(passwordHash),
			AvatarKey:    avatarKey,
		},
	)
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		UserID:      result.UserID,
		CompanyID:   result.CompanyID,
		DisplayName: displayName,
		Email:       email,
		CompanyName: companyName,
		Role:        "owner",
	}, nil
}
