package middleware

import (
	"be-education/config"
	"be-education/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtUtil *utils.JWTUtil
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	jwtUtil := utils.NewJWTUtil(cfg.SecretKey)
	return &AuthMiddleware{jwtUtil: jwtUtil}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format (Expected 'Bearer <token>')"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		claims, err := m.jwtUtil.ParseJWTToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		utils.SetUserClaimsToContext(c, claims)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserClaims, ok := utils.GetCurrentUserClaims(c)
		if !ok || currentUserClaims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found. Authentication required."})
			c.Abort()
			return
		}

		isAuthorized := false
		for _, role := range roles {
			if currentUserClaims.Role == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource. Insufficient role."})
			c.Abort()
			return
		}

		c.Next()
	}
}
