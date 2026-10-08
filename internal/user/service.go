package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/segoSambel/learn-golang-restapi/internal/database"
	db "github.com/segoSambel/learn-golang-restapi/internal/database/generated"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken = errors.New("email already registered")
)

type Store interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return UserResponse{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.store.CreateUser(
		ctx,
		db.CreateUserParams{
			Email:        email,
			PasswordHash: string(passwordHash),
		},
	)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return UserResponse{}, ErrEmailTaken
		}
		return UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
