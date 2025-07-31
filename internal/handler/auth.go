package handler

import (
	"chow/internal/model"
	"chow/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// validate data
	var req model.LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate data", Detail: err.Error()})
		return
	}

	// login user
	user, accessToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to login"})
		return
	}

	res := model.LoginUserRes{User: *user, AccessToken: accessToken}
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Login successful", Data: res})
}

func (h *AuthHandler) Register(c *gin.Context) {
	// validate data
	var req model.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate data", Detail: err.Error()})
		return
	}

	// register user
	user, err := h.authService.Register(c.Request.Context(), &model.User{Email: req.Email, Username: req.Username, Password: req.Password})
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExist) {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to register account"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Account registration successful", Data: user})
}
