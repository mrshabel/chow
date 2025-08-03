package handler

import (
	"chow/internal/model"
	"chow/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Middleware struct {
	AuthService *service.AuthService
}

func NewMiddleware(authService *service.AuthService) *Middleware {
	return &Middleware{AuthService: authService}
}

// Key type for context values
type contextKey string

const (
	// UserKey is the authenticated user info in the request context
	UserKey contextKey = "user"
)

// AuthMiddleware checks JWT tokens and adds user info to the request context
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token from Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Authorization header required"})
			return
		}

		// check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Invalid authorization format"})
			return
		}

		tokenString := parts[1]

		// validate the token
		claims, err := m.AuthService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: err.Error()})
			return
		}

		// Extract user ID, role, and username from claims
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Invalid user ID in token"})
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: service.ErrInvalidToken.Error()})
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Invalid role in token"})
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "Invalid username in token"})
			return
		}

		user := model.AuthenticatedUser{ID: userID, Username: username, Role: model.UserRole(role)}

		// Add user info to request context
		c.Set(string(UserKey), user)

		// process next handler
		c.Next()
	}
}

// GetCurrentUser retrieves the current user info from the request context
func GetCurrentUser(c *gin.Context) (model.AuthenticatedUser, bool) {
	user, ok := c.Get(string(UserKey))
	if !ok {
		return model.AuthenticatedUser{}, false
	}
	authUser, ok := user.(model.AuthenticatedUser)
	return authUser, ok
}

// GetCurrentAdmin retrieves the current admin info from the request context
func GetCurrentAdmin(c *gin.Context) (model.AuthenticatedUser, bool) {
	user, ok := c.Get(string(UserKey))
	if !ok {
		return model.AuthenticatedUser{}, false
	}
	authUser, ok := user.(model.AuthenticatedUser)
	if !ok {
		return model.AuthenticatedUser{}, false
	}
	if authUser.Role != model.Admin {
		return model.AuthenticatedUser{}, false
	}

	return authUser, ok
}
