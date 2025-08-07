package usecase

import (
	"Blog-API/internal/infrastructure/ai"
	"errors"
)

type AIUseCase struct {
	aiService *ai.AIService
}

func NewAIUseCase(aiService *ai.AIService) *AIUseCase {
	return &AIUseCase{
		aiService: aiService,
	}
}

// generate blog post
func (uc *AIUseCase) GenerateBlog(topic string) (string, error) {
	if topic == "" {
		return "", errors.New("topic is required")
	}

	if len(topic) < 3 {
		return "", errors.New("topic must be at least 3 characters long")
	}

	if len(topic) > 200 {
		return "", errors.New("topic is too long (max 200 characters)")
	}

	return uc.aiService.GenerateBlog(topic)
}

// enhance blog content
func (uc *AIUseCase) EnhanceBlog(blogContent string) (string, error) {
	if blogContent == "" {
		return "", errors.New("blog content is required")
	}

	if len(blogContent) < 50 {
		return "", errors.New("blog content must be at least 50 characters long")
	}

	if len(blogContent) > 10000 {
		return "", errors.New("blog content is too long (max 10,000 characters)")
	}

	return uc.aiService.EnhanceBlog(blogContent)
}

// suggest ideas based on keywords
func (uc *AIUseCase) SuggestBlogIdeas(keywords []string) (string, error) {
	if len(keywords) == 0 {
		return "", errors.New("at least one keyword is required")
	}

	if len(keywords) > 10 {
		return "", errors.New("too many keywords (max 10)")
	}

	for i, keyword := range keywords {
		if keyword == "" {
			return "", errors.New("keyword cannot be empty")
		}
		if len(keyword) < 2 {
			return "", errors.New("keyword must be at least 2 characters long")
		}
		if len(keyword) > 50 {
			return "", errors.New("keyword is too long (max 50 characters)")
		}
		// Remove duplicates by checking previous keywords
		for j := 0; j < i; j++ {
			if keywords[j] == keyword {
				return "", errors.New("duplicate keywords are not allowed")
			}
		}
	}

	return uc.aiService.SuggestBlogIdeas(keywords)
}
