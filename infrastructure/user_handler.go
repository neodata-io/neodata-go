package api

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rubenclaes/horse-rider-training/internal/user/app"
	"github.com/rubenclaes/horse-rider-training/internal/user/interfaces/dto"

	"github.com/rubenclaes/horse-rider-training/pkg/utils"
)

type HTTPHandler struct {
	UserService *app.UserService
}

// NewHTTPHandler creates a new HTTPHandler with UserService dependency.
func NewHTTPHandler(app *fiber.App, userService *app.UserService) *HTTPHandler {
	return &HTTPHandler{UserService: userService}
}

// HealthCheck is a simple handler for a health check endpoint.
func (h *HTTPHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("app is running")
}

func (h *HTTPHandler) CreateUser(c *fiber.Ctx) error {
	// Parse the request body into a User model
	var request dto.UserCreationDTO
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}

	// Call the Userapp to create a user
	user, err := h.UserService.CreateUser(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID handles HTTP requests to retrieve a user by their unique ID.
func (h *HTTPHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	user, err := h.UserService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// FetchUsers handles HTTP requests to fetch a list of users.
func (h *HTTPHandler) FetchUsers(c *fiber.Ctx) error {
	// Parse parameters
	searchQuery := c.Query("q", "")
	selectedStatus := c.Query("status", "")
	sortBy := c.Query("sortBy", "created_at")
	orderBy := c.Query("orderBy", "asc")
	itemsPerPage := utils.ParseQueryParam(c.Query("itemsPerPage"), 10)
	page := utils.ParseQueryParam(c.Query("page"), 1)

	users, totalCount, err := h.UserService.GetUsers(c.Context(), searchQuery, selectedStatus, itemsPerPage, page, sortBy, orderBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Return the list of users and total count in a JSON response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_count": totalCount,
		"users":       users,
	})
}

// UpdateUser handles HTTP requests to update an existing user's details.
func (h *HTTPHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var request dto.UserUpdateDTO
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Convert DTO to User entity
	user := request.ToEntity(userID)

	if err := h.UserService.UpdateUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DeleteUser handles HTTP requests to delete a user by their ID.
func (h *HTTPHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	if err := h.UserService.DeleteUser(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
