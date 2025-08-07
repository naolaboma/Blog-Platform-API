package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AIService struct {
	apiKey  string
	baseURL string
}

type AIRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type AIResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func NewAIService() *AIService {
	apiKey := os.Getenv("GROQ_API_KEY")
	// if apiKey == "" {
	// 	apiKey = ""
	// }

	return &AIService{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent",
	}
}

func (g *AIService) GenerateContent(prompt string) (string, error) {
	request := AIRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", g.baseURL, g.apiKey)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response AIResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w, body: %s", err, string(bodyBytes))
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated, response: %s", string(bodyBytes))
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// generate blog post
func (g *AIService) GenerateBlog(topic string) (string, error) {
	prompt := fmt.Sprintf(`Write a complete blog post about "%s". 
	Include:
	- An engaging title
	- Introduction
	- Main content with 3-4 sections
	- Conclusion
	- Make it informative, engaging, and well-structured.
	Format it with proper headings and paragraphs.`, topic)

	return g.GenerateContent(prompt)
}

// enhance blog content
func (g *AIService) EnhanceBlog(blogContent string) (string, error) {
	prompt := fmt.Sprintf(`Analyze this blog post and provide specific suggestions to improve it:

	%s

	Please provide:
	1. Content improvements (structure, flow, clarity)
	2. SEO suggestions
	3. Engagement tips
	4. Specific edits or additions
	5. Overall rating and areas of strength

	Be constructive and specific.`, blogContent)

	return g.GenerateContent(prompt)
}

// suggest ideas based on keywords
func (g *AIService) SuggestBlogIdeas(keywords []string) (string, error) {
	keywordStr := ""
	for i, keyword := range keywords {
		if i > 0 {
			keywordStr += ", "
		}
		keywordStr += keyword
	}

	prompt := fmt.Sprintf(`Generate 10 creative blog post ideas based on these keywords: %s

	For each idea, provide:
	1. A catchy title
	2. Brief description (2-3 sentences)
	3. Target audience
	4. Why it would be engaging

	Make the ideas diverse, practical, and engaging.`, keywordStr)

	return g.GenerateContent(prompt)
}
