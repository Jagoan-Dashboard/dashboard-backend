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
        if err == usecase.ErrInactiveUser {
            return response.Forbidden(c, "User account is inactive", err)
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




func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
    requesterID := c.Locals("userID").(string)

    result, err := h.authUseCase.GetAllUsers(c.Context(), requesterID)
    if err != nil {
        if err == usecase.ErrForbidden {
            return response.Forbidden(c, "Only superadmin can access this resource", err)
        }
        return response.InternalError(c, "Failed to get users", err)
    }

    return response.Success(c, "Users retrieved successfully", result)
}


func (h *AuthHandler) GetUserByID(c *fiber.Ctx) error {
    requesterID := c.Locals("userID").(string)
    targetUserID := c.Params("id")

    result, err := h.authUseCase.GetUserDetailByID(c.Context(), requesterID, targetUserID)
    if err != nil {
        if err == usecase.ErrForbidden {
            return response.Forbidden(c, "Only superadmin can access this resource", err)
        }
        if err == usecase.ErrUserNotFound {
            return response.NotFound(c, "User not found", err)
        }
        return response.InternalError(c, "Failed to get user", err)
    }

    return response.Success(c, "User retrieved successfully", result)
}


func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
    requesterID := c.Locals("userID").(string)

    var req dto.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    result, err := h.authUseCase.CreateUser(c.Context(), requesterID, &req)
    if err != nil {
        if err == usecase.ErrForbidden {
            return response.Forbidden(c, "Only superadmin can create users", err)
        }
        if err == usecase.ErrUserExists {
            return response.Conflict(c, "User already exists", err)
        }
        return response.InternalError(c, "Failed to create user", err)
    }

    return response.Created(c, "User created successfully", result)
}


func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
    requesterID := c.Locals("userID").(string)
    targetUserID := c.Params("id")

    var req dto.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    result, err := h.authUseCase.UpdateUserRole(c.Context(), requesterID, targetUserID, &req)
    if err != nil {
        if err == usecase.ErrForbidden {
            return response.Forbidden(c, "Only superadmin can update users", err)
        }
        if err == usecase.ErrUserNotFound {
            return response.NotFound(c, "User not found", err)
        }
        return response.InternalError(c, "Failed to update user", err)
    }

    return response.Success(c, "User updated successfully", result)
}


func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
    requesterID := c.Locals("userID").(string)
    targetUserID := c.Params("id")

    err := h.authUseCase.DeleteUser(c.Context(), requesterID, targetUserID)
    if err != nil {
        if err == usecase.ErrForbidden {
            return response.Forbidden(c, "Only superadmin can delete users", err)
        }
        if err == usecase.ErrUserNotFound {
            return response.NotFound(c, "User not found", err)
        }
        if err == usecase.ErrCannotDeleteSelf {
            return response.BadRequest(c, "Cannot delete your own account", err)
        }
        return response.InternalError(c, "Failed to delete user", err)
    }

    return response.Success(c, "User deleted successfully", nil)
}