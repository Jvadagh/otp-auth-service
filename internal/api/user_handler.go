package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jvadagh/otp-auth-service/internal/repository"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userRepo *repository.UserRepo
}

func NewUserHandler(userRepo *repository.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "user not found")
	}
	return c.JSON(user)
}

// ListUsers godoc
// @Summary List users
// @Description Retrieve a paginated list of users with optional search
// @Tags Users
// @Produce json
// @Param search query string false "Search by phone"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string "Database error"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	users, total, err := h.userRepo.List(search, page, limit)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "db error")
	}
	return c.JSON(fiber.Map{
		"items": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
