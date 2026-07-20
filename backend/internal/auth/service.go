package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"flowsync-pulse/backend/internal/token"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrJWTSecretNotSet    = errors.New("JWT_SECRET is not configured")
)

const accessTokenDuration = time.Hour

type AuthRepository interface {
	FindLoginUserByEmail(
		ctx context.Context,
		email string,
	) (LoginUserRecord, error)

	FindCurrentUser(
		ctx context.Context,
		userID uint64,
		companyID uint64,
	) (CurrentUserRecord, error)
}

type Service struct {
	repository AuthRepository
	jwtSecret  []byte
}

func NewService(repository AuthRepository) *Service {
	return &Service{
		repository: repository,
		jwtSecret:  []byte(os.Getenv("JWT_SECRET")),
	}
}

func (s *Service) Login(
	ctx context.Context,
	request LoginRequest,
) (LoginResponse, error) {
	email := strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	record, err := s.repository.FindLoginUserByEmail(
		ctx,
		email,
	)
	if err != nil {
		if errors.Is(err, ErrLoginUserNotFound) {
			return LoginResponse{}, ErrInvalidCredentials
		}

		return LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(record.PasswordHash),
		[]byte(request.Password),
	); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, expiresIn, err := s.generateAccessToken(record)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: LoginUser{
			ID:          record.UserID,
			DisplayName: record.DisplayName,
			Email:       record.Email,
			Role:        record.Role,
		},
		Company: LoginCompany{
			ID:   record.CompanyID,
			Name: record.CompanyName,
		},
	}, nil
}

func (s *Service) CurrentUser(
	ctx context.Context,
	userID uint64,
	companyID uint64,
) (CurrentUserResponse, error) {
	record, err := s.repository.FindCurrentUser(
		ctx,
		userID,
		companyID,
	)
	if err != nil {
		return CurrentUserResponse{}, err
	}

	return CurrentUserResponse{
		User: CurrentUserCompanyMember{
			ID:          record.UserID,
			DisplayName: record.DisplayName,
			Email:       record.Email,
			Role:        record.Role,
		},
		Company: CurrentUserCompany{
			ID:   record.CompanyID,
			Name: record.CompanyName,
		},
	}, nil
}

func (s *Service) generateAccessToken(
	record LoginUserRecord,
) (string, int64, error) {
	if len(s.jwtSecret) < 32 {
		return "", 0, ErrJWTSecretNotSet
	}

	now := time.Now()
	expiresAt := now.Add(accessTokenDuration)

	claims := token.AccessTokenClaims{
		CompanyID: record.CompanyID,
		Role:      record.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(record.UserID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "flowsync-pulse",
		},
	}

	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", 0, fmt.Errorf(
			"failed to sign access token: %w",
			err,
		)
	}

	return signedToken, int64(accessTokenDuration.Seconds()), nil
}
