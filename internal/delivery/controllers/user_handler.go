package controllers

import (
	"net/http"
	"strings"

	"Blog-API/internal/domain"
	"Blog-API/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// promotion , demotion and profile picture//
func (h *UserHandler) PromoteUser(c *gin.Context) {
	adminUserID, _ := middleware.GetUserIDFromContext(c)
	targetUserID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid target user ID"})
		return
	}
	err = h.userUseCase.UpdateRole(adminUserID, targetUserID, domain.RoleAdmin)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "forbidden") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, domain.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Promoted to admin successfully"})
}
func (h *UserHandler) DemoteUser(c *gin.Context) {
	adminUserID, _ := middleware.GetUserIDFromContext(c)
	targetUserID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}
	err = h.userUseCase.UpdateRole(adminUserID, targetUserID, domain.RoleUser)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "forbidden") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, domain.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin demoted to user successfully"})
}
func (h *UserHandler) UploadProfilePicture(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "User not authenticated"})
		return
	}
	file, handler, err := c.Request.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "File is required"})
		return
	}
	defer file.Close()

	const maxFileSize = 5 * 1024 * 1024
	if handler.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "file size exceeds the limit of 5MB"})
		return
	}
	contentType := handler.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "invalid file type. only JPG and PNG are alowed"})
		return
	}

	updatedUser, err := h.userUseCase.UploadProfilePicture(userID, file, handler)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedUser)

}

//promotion, demotion and profile picture//

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "token is required"})
		return
	}

	if err := h.userUseCase.VerifyEmail(token); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid or expired") {
			status = http.StatusBadRequest
		}
		c.JSON(status, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *UserHandler) SendVerificationEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	if err := h.userUseCase.SendVerificationEmail(req.Email); err != nil {
		// Return a generic response to avoid email enumeration
		if strings.Contains(err.Error(), "verification email has been sent") || strings.Contains(err.Error(), "already verified") {
			c.JSON(http.StatusOK, domain.EmailVerificationResponse{Message: "If an account exists for this email, a verification email has been sent."})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.EmailVerificationResponse{Message: "If an account exists for this email, a verification email has been sent."})
}

func (h *UserHandler) SendPasswordResetEmail(c *gin.Context) {
	var req domain.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	if err := h.userUseCase.SendPasswordResetEmail(req.Email); err != nil {
		// Generic response to avoid email enumeration
		if strings.Contains(err.Error(), "password reset email has been sent") {
			c.JSON(http.StatusOK, domain.PasswordResetResponse{Message: "If an account exists for this email, a password reset email has been sent."})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.PasswordResetResponse{Message: "If an account exists for this email, a password reset email has been sent."})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req domain.NewPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Invalid request data: " + err.Error()})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "Validation failed: " + err.Error()})
		return
	}

	if err := h.userUseCase.ResetPassword(req.Token, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid or expired") || strings.Contains(err.Error(), "password must") {
			status = http.StatusBadRequest
		}
		c.JSON(status, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.NewPasswordResponse{Message: "Password reset successful"})
}
