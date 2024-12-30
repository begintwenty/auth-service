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

type AuthRepo interface {
	FetchUserByID(ctx context.Context, userId domain.UserID) (*domain.Auth, error)
	Store(ctx context.Context, auth domain.Auth) error
}

type Service struct {
	authRepo AuthRepo
}

type User struct {
	UserId string
}

func New(authRepo AuthRepo) *Service {
	return &Service{
		authRepo: authRepo,
	}
}

func (s *Service) Store(ctx context.Context, auth domain.Auth) error {
	return s.authRepo.Store(ctx, auth)
}

func (s *Service) Authcheck(permissions ...domain.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerPrefix = "Bearer "
		authHeader := c.GetHeader("Authorization")

		// var tokenString string
		// if authHeader != "" && strings.HasPrefix(authHeader, bearerPrefix)  {
		// 	tokenString = strings.TrimPrefix(authHeader, bearerPrefix)
		// }else {
		// 	return nil, errors.New("bearer token missing or invalid")
		// }

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

		verifiedUserID, err := domain.VerifyUserID(unVerifiedUserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		auth, err := s.authRepo.FetchUserByID(c, *verifiedUserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		hasValidPerm := s.CheckPermissions(auth, permissions)
		if !hasValidPerm {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid permissions"})
			return
		}
		setCurrentAuth(c, *auth)
		c.Next()
	}
}

func setCurrentAuth(c *gin.Context, auth domain.Auth) {
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

func (a *Service) CheckPermissions(auth *domain.Auth, permissions []domain.Permission) bool {
	for _, permission := range permissions {
		if hasPermission, exists := auth.Permissions[permission]; !exists || !hasPermission {
			// Permission is either not present or explicitly set to false
			return false
		}
	}
	// All required permissions are granted
	return true
}
