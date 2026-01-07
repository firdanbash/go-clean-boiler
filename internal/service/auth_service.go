package service

import (
	"errors"

	"github.com/firdanbash/go-clean-boiler/internal/domain"
	"github.com/firdanbash/go-clean-boiler/internal/dto/request"
	"github.com/firdanbash/go-clean-boiler/internal/dto/response"
	"github.com/firdanbash/go-clean-boiler/internal/repository"
	"github.com/firdanbash/go-clean-boiler/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret, jwtExpiry string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register registers a new user
func (s *authService) Register(req *request.RegisterRequest) (*response.AuthResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		User: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}, nil
}

// Login authenticates a user and returns a token
func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		User: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}, nil
}

// generateToken generates a JWT token for the user
func (s *authService) generateToken(user *domain.User) (string, error) {
	// Parse JWT expiration duration
	duration, err := jwt.ParseDuration(s.jwtExpiry)
	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(user.ID, user.Email, s.jwtSecret, duration)
}
