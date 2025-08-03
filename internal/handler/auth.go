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

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginUserReq true "Login credentials"
// @Success 200 {object} model.SuccessResponse{data=model.LoginUserRes} "Login successful"
// @Failure 401 {object} model.ErrorResponse "Invalid credentials"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /auth/login [post]
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

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterUserReq true "User registration details"
// @Success 201 {object} model.SuccessResponse{data=model.User} "Account registration successful"
// @Failure 401 {object} model.ErrorResponse "User already exists"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /auth/register [post]
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

	c.JSON(http.StatusCreated, model.SuccessResponse{Message: "Account registration successful", Data: user})
}
