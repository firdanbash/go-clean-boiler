package handler

import (
	"github.com/firdanbash/go-clean-boiler/internal/dto/request"
	"github.com/firdanbash/go-clean-boiler/internal/service"
	"github.com/firdanbash/go-clean-boiler/pkg/response"
	"github.com/firdanbash/go-clean-boiler/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Registration request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	if !validator.BindAndValidate(c, &req) {
		return
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User registered successfully", result)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if !validator.BindAndValidate(c, &req) {
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Login successful", result)
}
