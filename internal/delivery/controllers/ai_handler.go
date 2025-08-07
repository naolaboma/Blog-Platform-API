package controllers

import (
	"Blog-API/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AIHandler struct {
	aiUseCase *usecase.AIUseCase
	validate  *validator.Validate
}

func NewAIHandler(aiUseCase *usecase.AIUseCase) *AIHandler {
	return &AIHandler{
		aiUseCase: aiUseCase,
		validate:  validator.New(),
	}
}

type GenerateBlogRequest struct {
	Topic string `json:"topic" validate:"required,min=3,max=200"`
}

type EnhanceBlogRequest struct {
	Content string `json:"content" validate:"required,min=50,max=10000"`
}

type SuggestBlogIdeasRequest struct {
	Keywords []string `json:"keywords" validate:"required,min=1,max=10,dive,min=2,max=50"`
}

func (h *AIHandler) GenerateBlog(c *gin.Context) {
	var req GenerateBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	generatedBlog, err := h.aiUseCase.GenerateBlog(req.Topic)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "topic is required" ||
			err.Error() == "topic must be at least 3 characters long" ||
			err.Error() == "topic is too long (max 200 characters)" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Blog generated successfully",
		"topic":   req.Topic,
		"content": generatedBlog,
	})
}

func (h *AIHandler) EnhanceBlog(c *gin.Context) {
	var req EnhanceBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	suggestions, err := h.aiUseCase.EnhanceBlog(req.Content)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "blog content is required" ||
			err.Error() == "blog content must be at least 50 characters long" ||
			err.Error() == "blog content is too long (max 10,000 characters)" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Blog enhancement suggestions generated",
		"suggestions": suggestions,
	})
}

func (h *AIHandler) SuggestBlogIdeas(c *gin.Context) {
	var req SuggestBlogIdeasRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	ideas, err := h.aiUseCase.SuggestBlogIdeas(req.Keywords)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "at least one keyword is required" ||
			err.Error() == "too many keywords (max 10)" ||
			err.Error() == "keyword cannot be empty" ||
			err.Error() == "keyword must be at least 2 characters long" ||
			err.Error() == "keyword is too long (max 50 characters)" ||
			err.Error() == "duplicate keywords are not allowed" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Blog ideas generated successfully",
		"keywords": req.Keywords,
		"ideas":    ideas,
	})
}
