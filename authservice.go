package authservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dev-mantas/authservice/domain"
	"github.com/dev-mantas/authservice/pkg/token"

	"github.com/gin-gonic/gin"
)

type Authenticable interface {
	GetUserID() string
	HasPermission(permission string) bool
}

type UserRepo[T Authenticable] interface {
	FetchUserByIDAsString(ctx context.Context, userID string) (*T, error)
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

		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
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

		// Check permissions via the interface
		// for _, perm := range permissions {
		// 	if !user.HasPermission(perm) {
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid permissions"})
		// 		return
		// 	}
		// }

		c.Set("currentUser", user)
		c.Next()
	}
}

func (a *Service[T]) GetUserFromContext(c *gin.Context) *domain.Auth {
	currentAuthInterface, exists := c.Get("currentUser")
	if !exists {
		fmt.Println(" Current user doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return nil
	}
	currentAuth, ok := currentAuthInterface.(*domain.Auth)
	if !ok {
		fmt.Println("Couldn't user current user")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return nil
	}

	return currentAuth
}
