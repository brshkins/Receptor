package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/model"
	"receptor/backend/internal/repository"
)

const jwtTTL = 24 * time.Hour

type authService struct {
	users      repository.UsersRepository
	jwtSecret  []byte
}

func NewAuthService(users repository.UsersRepository, jwtSecret string) AuthService {
	return &authService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Register(ctx context.Context, in dto.RegisterInput) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || strings.TrimSpace(in.Password) == "" {
		return nil, ErrInvalidCredentials
	}

	existing, err := s.users.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         strings.TrimSpace(in.Name),
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}

	token, err := s.signJWT(u.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token}, nil
}

func (s *authService) Login(ctx context.Context, in dto.LoginInput) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.signJWT(u.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token}, nil
}

func (s *authService) Me(ctx context.Context, userID int64) (*dto.MeResponse, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.MeResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, nil
}

func (s *authService) signJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(jwtTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}
