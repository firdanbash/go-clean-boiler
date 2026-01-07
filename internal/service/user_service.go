package service

import (
	"errors"

	"github.com/firdanbash/go-clean-boiler/internal/domain"
	"github.com/firdanbash/go-clean-boiler/internal/dto/request"
	"github.com/firdanbash/go-clean-boiler/internal/dto/response"
	"github.com/firdanbash/go-clean-boiler/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(req *request.CreateUserRequest) (*response.UserResponse, error)
	GetByID(id uint) (*response.UserResponse, error)
	GetAll(page, perPage int) ([]response.UserResponse, int64, error)
	Update(id uint, req *request.UpdateUserRequest) (*response.UserResponse, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create creates a new user
func (s *userService) Create(req *request.CreateUserRequest) (*response.UserResponse, error) {
	// Check if email already exists
	_, err := s.repo.FindByEmail(req.Email)
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

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// GetByID gets a user by ID
func (s *userService) GetByID(id uint) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// GetAll gets all users with pagination
func (s *userService) GetAll(page, perPage int) ([]response.UserResponse, int64, error) {
	offset := (page - 1) * perPage
	users, total, err := s.repo.FindAll(perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(&user)
	}

	return userResponses, total, nil
}

// Update updates a user
func (s *userService) Update(id uint, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, err := s.repo.FindByEmail(req.Email)
		if err == nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// Delete deletes a user
func (s *userService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

// toUserResponse converts domain.User to response.UserResponse
func (s *userService) toUserResponse(user *domain.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
