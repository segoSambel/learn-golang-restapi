package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segoSambel/learn-golang-restapi/internal/database"
	db "github.com/segoSambel/learn-golang-restapi/internal/database/generated"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken = errors.New("email already registered")
	ErrNotFound   = errors.New("user not found")
)

type Store interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error)
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
	if database.IsUniqueViolation(err) {
		return UserResponse{}, ErrEmailTaken
	}
	if err != nil {
		return UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return toResponse(user), nil
}

func (s *Service) Get(ctx context.Context, id string) (UserResponse, error) {
	var userID pgtype.UUID
	if err := userID.Scan(id); err != nil {
		return UserResponse{}, ErrNotFound
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return UserResponse{}, ErrNotFound
	}
	if err != nil {
		return UserResponse{}, fmt.Errorf("get user: %w", err)
	}

	return toResponse(user), nil
}

func toResponse(user db.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}
}
