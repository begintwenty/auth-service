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
		var tokenString string

		if authHeader := c.GetHeader("Authorization"); strings.HasPrefix(authHeader, bearerPrefix) {
			tokenString = strings.TrimPrefix(authHeader, bearerPrefix)
		} else if qToken := c.Query("token"); qToken != "" {
			tokenString = qToken
		} else {
			host := c.Request.Host
			subdomain := strings.Split(host, ".")[0]
			fmt.Println(host)
			fmt.Println(subdomain)
			switch subdomain {
			case "admin":
				tokenString, _ = c.Cookie("AX-JWT")
			case "app":
				tokenString, _ = c.Cookie("X-JWT")
			default:
				if userToken, err := c.Cookie("X-JWT"); err == nil {
					tokenString = userToken
				} else if adminToken, err := c.Cookie("AX-JWT"); err == nil {
					tokenString = adminToken
				}
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		userID, err := token.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		user, err := s.userRepo.FetchUserByIDAsString(c, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
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
