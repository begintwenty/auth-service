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

type UserRepo interface {
	FetchUserByUserIDAsString(ctx context.Context, userId string) (*User, error)
	Store(ctx context.Context, auth domain.Auth) error
}

type Service struct {
	userRepo UserRepo
}

type User struct {
	UserId      string
	Permissions map[string]bool
}

func New(userRepo UserRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) Authcheck(permissions ...string) gin.HandlerFunc {
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

		user, err := s.userRepo.FetchUserByID(c, unVerifiedUserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		hasValidPerm := s.CheckPermissions(user, permissions)
		if !hasValidPerm {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid permissions"})
			return
		}
		setCurrentAuth(c, *user)
		c.Next()
	}
}

func setCurrentAuth(c *gin.Context, auth User) {
	// currentUser := user.ToCurrentUser()
	c.Set("currentAuth", &auth)

}

func (a *Service) GetAuthFromContext(c *gin.Context) *domain.Auth {
	currentAuthInterface, exists := c.Get("currentAuth")
	if !exists {
		fmt.Println(" Current auth doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return nil
	}
	currentAuth, ok := currentAuthInterface.(*domain.Auth)
	if !ok {
		fmt.Println("Couldn't auth current user")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return nil
	}

	return currentAuth
}

func (a *Service) CheckPermissions(user *User, permissions []string) bool {
	for _, permission := range permissions {
		if hasPermission, exists := user.Permissions[permission]; !exists || !hasPermission {
			// Permission is either not present or explicitly set to false
			return false
		}
	}
	// All required permissions are granted
	return true
}
