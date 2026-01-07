package handler

import (
	"strconv"

	"github.com/firdanbash/go-clean-boiler/internal/dto/request"
	"github.com/firdanbash/go-clean-boiler/internal/service"
	"github.com/firdanbash/go-clean-boiler/pkg/response"
	"github.com/firdanbash/go-clean-boiler/pkg/validator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "Create user request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest

	if !validator.BindAndValidate(c, &req) {
		return
	}

	result, err := h.userService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User created successfully", result)
}

// GetAll godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, total, err := h.userService.GetAll(page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	pagination := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
	}

	response.Paginated(c, "Users retrieved successfully", users, pagination)
}

// GetByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

// Update godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body request.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req request.UpdateUserRequest
	if !validator.BindAndValidate(c, &req) {
		return
	}

	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User updated successfully", user)
}

// Delete godoc
// @Summary Delete user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User deleted successfully", nil)
}
