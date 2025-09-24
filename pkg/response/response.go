package response

import (
	"math"

	"building-report-backend/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains metadata for paginated responses
type Meta struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	response := APIResponse{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	return c.JSON(response)
}

// SuccessWithMeta sends a successful response with pagination metadata
func SuccessWithMeta(c *fiber.Ctx, data interface{}, meta *Meta, message ...string) error {
	response := APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	return c.JSON(response)
}

// Error sends an error response
func Error(c *fiber.Ctx, err error) error {
	var appError *errors.AppError

	if errors.IsAppError(err) {
		appError = errors.GetAppError(err)
	} else {
		appError = errors.GetAppError(err)
	}

	response := APIResponse{
		Success: false,
		Message: appError.Message,
		Error: map[string]interface{}{
			"code":    appError.Code,
			"message": appError.Message,
		},
	}

	if appError.Details != "" {
		response.Error = map[string]interface{}{
			"code":    appError.Code,
			"message": appError.Message,
			"details": appError.Details,
		}
	}

	return c.Status(appError.HTTPStatus).JSON(response)
}

// ValidationError sends a validation error response
func ValidationError(c *fiber.Ctx, message string, details ...string) error {
	response := APIResponse{
		Success: false,
		Message: message,
		Error: map[string]interface{}{
			"code":    errors.ErrCodeValidationFailed,
			"message": message,
		},
	}

	if len(details) > 0 {
		response.Error = map[string]interface{}{
			"code":    errors.ErrCodeValidationFailed,
			"message": message,
			"details": details[0],
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// NotFound sends a not found error response
func NotFound(c *fiber.Ctx, resource string) error {
	appError := errors.NewNotFoundError(resource)
	return Error(c, appError)
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *fiber.Ctx, message ...string) error {
	msg := "Authentication required"
	if len(message) > 0 {
		msg = message[0]
	}

	response := APIResponse{
		Success: false,
		Message: msg,
		Error: map[string]interface{}{
			"code":    errors.ErrCodeUnauthorized,
			"message": msg,
		},
	}

	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

// Forbidden sends a forbidden error response
func Forbidden(c *fiber.Ctx, message ...string) error {
	msg := "Access forbidden"
	if len(message) > 0 {
		msg = message[0]
	}

	response := APIResponse{
		Success: false,
		Message: msg,
		Error: map[string]interface{}{
			"code":    errors.ErrCodeForbidden,
			"message": msg,
		},
	}

	return c.Status(fiber.StatusForbidden).JSON(response)
}

// InternalError sends an internal server error response
func InternalError(c *fiber.Ctx, err error) error {
	response := APIResponse{
		Success: false,
		Message: "Internal server error",
		Error: map[string]interface{}{
			"code":    errors.ErrCodeInternalError,
			"message": "Internal server error",
		},
	}

	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

// NewMeta creates pagination metadata
func NewMeta(page, pageSize, totalRecords int) *Meta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return &Meta{
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

// Created sends a resource created response
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	response := APIResponse{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	} else {
		response.Message = "Resource created successfully"
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Updated sends a resource updated response
func Updated(c *fiber.Ctx, data interface{}, message ...string) error {
	response := APIResponse{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	} else {
		response.Message = "Resource updated successfully"
	}

	return c.JSON(response)
}

// Deleted sends a resource deleted response
func Deleted(c *fiber.Ctx, message ...string) error {
	response := APIResponse{
		Success: true,
	}

	if len(message) > 0 {
		response.Message = message[0]
	} else {
		response.Message = "Resource deleted successfully"
	}

	return c.JSON(response)
}