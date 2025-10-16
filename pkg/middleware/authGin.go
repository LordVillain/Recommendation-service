// internal/middleware/gin_auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/LordVillain/Recommendation-service/configs"
	"github.com/ShopOnGO/ShopOnGO/pkg/jwt"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/gin-gonic/gin"
)

type key string
const (
	ContextUserIDKey key = "ContextUserIDKey"
	ContextRolesKey key = "ContextRolesKey"
)

func GinAuthMiddleware(config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("❌ No valid Bearer prefix")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data, err := jwt.NewJWT(config.OAuth.Secret).Parse(token)

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				logger.Error("❌ Token expired:", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
			logger.Error("❌ Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !isValid {
			logger.Error("❌ Token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		logger.Info("✅ Token is valid for:", data.UserID)

		// Кладём в контекст Gin (а не в http.Request.Context)
		c.Set(string(ContextUserIDKey), data.UserID)
		c.Set(string(ContextRolesKey), data.Role)

		// Продолжаем обработку
		c.Next()
	}
}