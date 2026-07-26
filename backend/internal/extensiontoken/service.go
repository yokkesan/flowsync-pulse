package extensiontoken

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	extensionTokenByteLength = 32
	extensionTokenValidDays  = 90
)

var ErrInvalidExtensionTokenName = errors.New(
	"invalid extension token name",
)

type repositoryInterface interface {
	Create(
		ctx context.Context,
		params CreateParams,
	) (CreateResult, error)

	FindAllByUser(
		ctx context.Context,
		userID uint64,
		companyID uint64,
	) ([]Response, error)

	Revoke(
		ctx context.Context,
		tokenID uint64,
		userID uint64,
		companyID uint64,
	) error
}

type Service struct {
	repository repositoryInterface
}

func NewService(
	repository repositoryInterface,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	userID uint64,
	companyID uint64,
	request CreateRequest,
) (CreateResponse, error) {
	name := strings.TrimSpace(request.Name)

	if name == "" || len(name) > 100 {
		return CreateResponse{},
			ErrInvalidExtensionTokenName
	}

	plainToken, err := generateExtensionToken()
	if err != nil {
		return CreateResponse{}, err
	}

	tokenHash := hashExtensionToken(
		plainToken,
	)

	createdAt := time.Now().UTC()
	expiresAt := createdAt.AddDate(
		0,
		0,
		extensionTokenValidDays,
	)

	result, err := s.repository.Create(
		ctx,
		CreateParams{
			UserID:    userID,
			CompanyID: companyID,
			Name:      name,
			TokenHash: tokenHash,
			ExpiresAt: expiresAt,
		},
	)
	if err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{
		TokenID:   result.TokenID,
		Name:      name,
		Token:     plainToken,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}, nil
}

func (s *Service) List(
	ctx context.Context,
	userID uint64,
	companyID uint64,
) (ListResponse, error) {
	tokens, err := s.repository.FindAllByUser(
		ctx,
		userID,
		companyID,
	)
	if err != nil {
		return ListResponse{}, err
	}

	return ListResponse{
		Tokens: tokens,
	}, nil
}

func (s *Service) Revoke(
	ctx context.Context,
	tokenID uint64,
	userID uint64,
	companyID uint64,
) error {
	return s.repository.Revoke(
		ctx,
		tokenID,
		userID,
		companyID,
	)
}

func generateExtensionToken() (
	string,
	error,
) {
	randomBytes := make(
		[]byte,
		extensionTokenByteLength,
	)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf(
			"failed to generate extension token: %w",
			err,
		)
	}

	return base64.RawURLEncoding.EncodeToString(
		randomBytes,
	), nil
}

func hashExtensionToken(
	plainToken string,
) string {
	hashedToken := sha256.Sum256(
		[]byte(plainToken),
	)

	return hex.EncodeToString(
		hashedToken[:],
	)
}
