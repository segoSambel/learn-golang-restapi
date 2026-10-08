package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	db "github.com/segoSambel/learn-golang-restapi/internal/database/generated"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
)

type Store interface {
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}

type Service struct {
	store    Store
	secret   []byte
	tokenTTL time.Duration
}

func NewService(store Store, secret string, tokenTTL time.Duration) *Service {
	return &Service{
		store:    store,
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (TokenResponse, error) {
	user, err := s.store.GetUserByEmail(ctx, strings.TrimSpace(req.Email))
	if errors.Is(err, pgx.ErrNoRows) {
		return TokenResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return TokenResponse{}, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return TokenResponse{}, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(s.tokenTTL)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	accessToken, err := token.SignedString(s.secret)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("sign token: %w", err)
	}

	return TokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *Service) Verify(accessToken string) (string, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(
		accessToken,
		&claims,
		func(*jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || claims.Subject == "" {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}
