package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vietgs03/translate/backend/internal/config"
	"github.com/vietgs03/translate/backend/internal/errors"
	"github.com/vietgs03/translate/backend/internal/model"
	"github.com/vietgs03/translate/backend/internal/repository"
	"github.com/vietgs03/translate/backend/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*model.User, error)
	Login(ctx context.Context, input LoginInput) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	UpdateRole(ctx context.Context, userID uint, role string) (*model.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtConfig config.JWTConfig) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*model.User, error) {
	// Check if username exists
	if _, err := s.userRepo.GetByUsername(ctx, input.Username); err == nil {
		return nil, errors.NewValidationError("username already exists")
	}

	// Check if email exists
	if _, err := s.userRepo.GetByEmail(ctx, input.Email); err == nil {
		return nil, errors.NewValidationError("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "translator", // Change default role to translator
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("failed to create user: %v", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return "", errors.NewValidationError("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.NewValidationError("invalid credentials")
	}

	// Generate JWT token
	claims := &types.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtConfig.ExpiresIn) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.SecretKey), nil
	})
}

func (s *authService) UpdateRole(ctx context.Context, userID uint, role string) (*model.User, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Update role
	user.Role = role
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("failed to update user role: %v", err)
	}

	return user, nil
} 