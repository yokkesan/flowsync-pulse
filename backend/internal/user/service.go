package user

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordMismatch = errors.New(
	"password confirmation does not match",
)

type UserRepository interface {
	Register(
		ctx context.Context,
		params RegisterParams,
	) (RegisterResult, error)

	ListByCompanyID(
		ctx context.Context,
		companyID uint64,
	) ([]ListMemberResult, error)
}

type Service struct {
	repository UserRepository
}

func NewService(repository UserRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Register(
	ctx context.Context,
	companyID uint64,
	request RegisterRequest,
) (RegisterResponse, error) {
	displayName := strings.TrimSpace(
		request.DisplayName,
	)

	email := strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	if request.Password != request.PasswordConfirm {
		return RegisterResponse{}, ErrPasswordMismatch
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return RegisterResponse{}, err
	}

	result, err := s.repository.Register(
		ctx,
		RegisterParams{
			CompanyID:    companyID,
			DisplayName:  displayName,
			Email:        email,
			PasswordHash: string(passwordHash),
		},
	)
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		UserID:      result.UserID,
		CompanyID:   companyID,
		DisplayName: displayName,
		Email:       email,
		Role:        "owner",
	}, nil
}

func (s *Service) ListByCompanyID(
	ctx context.Context,
	companyID uint64,
) (ListResponse, error) {
	members, err := s.repository.ListByCompanyID(
		ctx,
		companyID,
	)
	if err != nil {
		return ListResponse{}, err
	}

	users := make(
		[]ListMemberResponse,
		0,
		len(members),
	)

	for _, member := range members {
		users = append(
			users,
			ListMemberResponse{
				UserID:      member.UserID,
				DisplayName: member.DisplayName,
				Email:       member.Email,
				Role:        member.Role,
				Status:      member.Status,
			},
		)
	}

	return ListResponse{
		Users: users,
	}, nil
}
