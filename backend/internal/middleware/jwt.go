package middleware

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userIDKey = "jwt_user_id"

// UserID returns the authenticated user id set by JWT middleware.
func UserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

// JWT validates Authorization: Bearer <token>, HS256 with secret, sets user id in context.
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(h, prefix) {
			respondUnauthorized(c)
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(h, prefix))
		if raw == "" {
			respondUnauthorized(c)
			return
		}
		token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			respondUnauthorized(c)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondUnauthorized(c)
			return
		}
		sub, _ := claims["sub"].(string)
		uid, err := strconv.ParseInt(sub, 10, 64)
		if err != nil || uid <= 0 {
			respondUnauthorized(c)
			return
		}
		c.Set(userIDKey, uid)
		c.Next()
	}
}

func respondUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
}
