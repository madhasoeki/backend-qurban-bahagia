package middleware

import (
	"net/http"
	"qurban/utils"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authorizedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Membutuhkan autentikasi"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluwarsa"})
			c.Abort()
			return
		}

		if len(authorizedRoles) > 0 && !slices.Contains(authorizedRoles, claims.Role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin akses"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}