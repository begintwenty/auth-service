package authservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dev-mantas/authservice/pkg/token"

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

		// If any permissions were passed in, check them
		for _, perm := range permissions {
			if !user.HasPermission(perm) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("Missing permission: %s", perm),
				})
				return
			}
		}

		// If all checks pass, set the current user in the context
		c.Set("currentUser", user)
		c.Next()
	}
}

func (a *Service[T]) GetUserFromContext(c *gin.Context) *T {
	currentUserInterface, exists := c.Get("currentUser")
	if !exists {
		fmt.Println("Current user doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return nil
	}

	currentUser, ok := currentUserInterface.(*T)
	if !ok {
		fmt.Println("Couldn't cast current user to *T")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return nil
	}

	return currentUser
}
