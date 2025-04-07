package authservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/begintwenty/auth-service/pkg/token"
	"github.com/gin-gonic/gin"
)

type Authenticable interface {
	GetUserID() string
	HasPermission(permission string) bool
}

type UserRepo[T Authenticable] interface {
	FetchUserByIDAsString(ctx context.Context, userID string) (T, error)
}

type Service[T Authenticable] struct {
	userRepo UserRepo[T]
}

func New[T Authenticable](userRepo UserRepo[T]) *Service[T] {
	return &Service[T]{
		userRepo: userRepo,
	}
}

func (s *Service[T]) Authcheck(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerPrefix = "Bearer "
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		switch {
		case strings.HasPrefix(authHeader, bearerPrefix):
			tokenString = strings.TrimPrefix(authHeader, bearerPrefix)

		case c.Query("token") != "":
			tokenString = c.Query("token")

		default:
			cookieToken, err := c.Cookie("X-JWT")
			if err != nil || cookieToken == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid token"})
				return
			}
			tokenString = cookieToken
		}

		unVerifiedUserID, err := token.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		user, err := s.userRepo.FetchUserByIDAsString(c, unVerifiedUserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		for _, perm := range permissions {
			if !user.HasPermission(perm) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("Missing permission: %s", perm),
				})
				return
			}
		}

		c.Set("currentUser", user)
		c.Next()
	}
}

func (a *Service[T]) GetUserFromContext(c *gin.Context) T {
	var zeroValue T

	currentUserInterface, exists := c.Get("currentUser")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return zeroValue
	}

	currentUser, ok := currentUserInterface.(T)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return zeroValue
	}

	return currentUser
}
