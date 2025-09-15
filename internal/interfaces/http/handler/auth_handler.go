
package handler

import (
    "building-report-backend/internal/application/dto"
    "building-report-backend/internal/application/usecase"
    "building-report-backend/internal/interfaces/response"
    
    "github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
    authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
    return &AuthHandler{
        authUseCase: authUseCase,
    }
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req dto.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    result, err := h.authUseCase.Register(c.Context(), &req)
    if err != nil {
        if err == usecase.ErrUserExists {
            return response.Conflict(c, "User already exists", err)
        }
        return response.InternalError(c, "Failed to register user", err)
    }

    return response.Success(c, "User registered successfully", result)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req dto.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    result, err := h.authUseCase.Login(c.Context(), &req)
    if err != nil {
        if err == usecase.ErrInvalidCredentials {
            return response.Unauthorized(c, "Invalid credentials", err)
        }
        return response.InternalError(c, "Failed to login", err)
    }

    return response.Success(c, "Login successful", result)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    
    user, err := h.authUseCase.GetUserByID(c.Context(), userID)
    if err != nil {
        return response.NotFound(c, "User not found", err)
    }

    return response.Success(c, "User profile retrieved", user)
}