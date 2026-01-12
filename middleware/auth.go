package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"worknote-api/contract"
	"worknote-api/services/auth_service"
	"worknote-api/utils/render"
)

const (
	// UserInfoKey is the context key for user info
	UserInfoKey = "user_info"
)

// AuthMiddleware validates JWE tokens and injects user info into context
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return render.Unauthorized(c, "missing authorization header")
	}

	// Extract Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return render.Unauthorized(c, "invalid authorization header format")
	}

	token := parts[1]

	// Decrypt and validate token
	claims, err := auth_service.DecryptJWEToken(token)
	if err != nil {
		return render.Unauthorized(c, "invalid or expired token")
	}

	// Inject user info into context
	userInfo := &contract.UserInfo{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}
	c.Locals(UserInfoKey, userInfo)

	return c.Next()
}

// GetUserFromContext retrieves user info from fiber context
func GetUserFromContext(c *fiber.Ctx) *contract.UserInfo {
	userInfo, ok := c.Locals(UserInfoKey).(*contract.UserInfo)
	if !ok {
		return nil
	}
	return userInfo
}
