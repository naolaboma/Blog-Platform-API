package controllers

import (
	"net/http"

	"Blog-API/internal/domain"
	"Blog-API/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
	validate    *validator.Validate
}

func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validate:    validator.New(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	user, err := h.userUseCase.Register(req.Username, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" || err.Error() == "user with this username already exists" {
			status = http.StatusConflict
		} else if err.Error() == "password must be at least 6 characters long" ||
			err.Error() == "password must contain at least one uppercase letter" ||
			err.Error() == "password must contain at least one lowercase letter" ||
			err.Error() == "password must contain at least one number" ||
			err.Error() == "password must contain at least one special character" {
			status = http.StatusBadRequest
		}
		c.JSON(status, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	response, err := h.userUseCase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "User not authenticated"})
		return
	}

	user, err := h.userUseCase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	response, err := h.userUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "User not authenticated"})
		return
	}

	err := h.userUseCase.Logout(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.LogoutResponse{
		Message: "Successfully logged out",
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	updatedUser, err := h.userUseCase.UpdateProfile(userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			status = http.StatusConflict
		} else if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, domain.ErrorResponse{Error: "email verification endpoint not implemented yet"})
}

func (h *UserHandler) SendVerificationEmail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, domain.ErrorResponse{Error: "send verification email endpoint not implemented yet"})
}

func (h *UserHandler) SendPasswordResetEmail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, domain.ErrorResponse{Error: "send password reset email endpoint not implemented yet"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, domain.ErrorResponse{Error: "password reset endpoint not implemented yet"})
}
