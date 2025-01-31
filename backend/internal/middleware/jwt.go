package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vietgs03/translate/backend/internal/errors"
	"github.com/vietgs03/translate/backend/internal/types"
)

func JWTAuth(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.NewUnauthorizedError("missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return errors.NewUnauthorizedError("invalid token format")
		}

		token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			return errors.NewUnauthorizedError("invalid token: %v", err)
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok || !token.Valid {
			return errors.NewUnauthorizedError("invalid token claims")
		}

		// Add claims to context
		c.Locals("user", claims)
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*types.JWTClaims)
		if !ok {
			return errors.NewUnauthorizedError("user not authenticated")
		}

		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return errors.NewUnauthorizedError("insufficient permissions")
	}
} 