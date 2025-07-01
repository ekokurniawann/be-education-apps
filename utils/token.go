package utils

import (
	"be-education/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secretKey []byte
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTUtil) GenerateJWTToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "be-education-app",
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			Audience:  []string{"users"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (j *JWTUtil) ParseJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

type ContextKey string

const (
	UserClaimsContextKey ContextKey = "userClaims"
)

func SetUserClaimsToContext(c *gin.Context, claims *Claims) {
	c.Set(string(UserClaimsContextKey), claims)
}

func GetCurrentUserClaims(c *gin.Context) (*Claims, bool) {
	val, ok := c.Get(string(UserClaimsContextKey))
	if !ok {
		return nil, false
	}
	claims, ok := val.(*Claims)
	return claims, ok
}
