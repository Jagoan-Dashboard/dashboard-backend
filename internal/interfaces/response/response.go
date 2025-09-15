
package response

import (
    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func BadRequest(c *fiber.Ctx, message string, err error) error {
    return c.Status(fiber.StatusBadRequest).JSON(Response{
        Success: false,
        Message: message,
        Error:   getErrorMessage(err),
    })
}

func Unauthorized(c *fiber.Ctx, message string, err error) error {
    return c.Status(fiber.StatusUnauthorized).JSON(Response{
        Success: false,
        Message: message,
        Error:   getErrorMessage(err),
    })
}

func Forbidden(c *fiber.Ctx, message string, err error) error {
    return c.Status(fiber.StatusForbidden).JSON(Response{
        Success: false,
        Message: message,
        Error:   getErrorMessage(err),
    })
}

func NotFound(c *fiber.Ctx, message string, err error) error {
    return c.Status(fiber.StatusNotFound).JSON(Response{
        Success: false,
        Message: message,
        Error:   getErrorMessage(err),
    })
}

func Conflict(c *fiber.Ctx, message string, err error) error {
    return c.Status(fiber.StatusConflict).JSON(Response{
        Success: false,
        Message: message,
        Error:   getErrorMessage(err),
    })
}

func InternalError(c *fiber.Ctx, message string, err error) error {
    return c.Status(fiber.StatusInternalServerError).JSON(Response{
        Success: false,
        Message: message,
        Error:   getErrorMessage(err),
    })
}

func ValidationError(c *fiber.Ctx, err error) error {
    var validationErrors []string
    
    if errs, ok := err.(validator.ValidationErrors); ok {
        for _, e := range errs {
            validationErrors = append(validationErrors, formatValidationError(e))
        }
    }
    
    return c.Status(fiber.StatusBadRequest).JSON(Response{
        Success: false,
        Message: "Validation failed",
        Error:   validationErrors,
    })
}

func formatValidationError(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return err.Field() + " is required"
    case "min":
        return err.Field() + " must be at least " + err.Param()
    case "max":
        return err.Field() + " must be at most " + err.Param()
    case "email":
        return err.Field() + " must be a valid email"
    default:
        return err.Field() + " is invalid"
    }
}

func getErrorMessage(err error) interface{} {
    if err != nil {
        return err.Error()
    }
    return nil
}