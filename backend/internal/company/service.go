package company

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidSlug = errors.New(
	"company slug must contain only lowercase letters, numbers, and hyphens",
)

var companySlugPattern = regexp.MustCompile(
	`^[a-z0-9]+(?:-[a-z0-9]+)*$`,
)

type CompanyCreator interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (CreateResult, error)
}

type Service struct {
	repository CompanyCreator
}

func NewService(repository CompanyCreator) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	request CreateRequest,
) (CreateResponse, error) {
	name := strings.TrimSpace(request.Name)
	slug := strings.ToLower(
		strings.TrimSpace(request.Slug),
	)

	if !companySlugPattern.MatchString(slug) {
		return CreateResponse{}, ErrInvalidSlug
	}

	result, err := s.repository.Create(
		ctx,
		CreateParams{
			Name: name,
			Slug: slug,
		},
	)
	if err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{
		CompanyID: result.CompanyID,
		Name:      name,
		Slug:      slug,
	}, nil
}
