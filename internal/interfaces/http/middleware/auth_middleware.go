// internal/interfaces/http/middleware/auth_middleware.go
package middleware

import (
    "strings"
    
    "building-report-backend/internal/infrastructure/auth"
    "building-report-backend/internal/interfaces/response"
    "github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtService auth.JWTService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return response.Unauthorized(c, "Missing authorization header", nil)
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        if tokenString == "" {
            return response.Unauthorized(c, "Invalid token format", nil)
        }

        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            return response.Unauthorized(c, "Invalid or expired token", err)
        }

        // Set user info in context
        c.Locals("userID", claims.UserID)
        c.Locals("username", claims.Username)
        c.Locals("role", claims.Role)

        return c.Next()
    }
}

func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("role").(string)
        
        for _, role := range roles {
            if userRole == role {
                return c.Next()
            }
        }
        
        return response.Forbidden(c, "Insufficient permissions", nil)
    }
}