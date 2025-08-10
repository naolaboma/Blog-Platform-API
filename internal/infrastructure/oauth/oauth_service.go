package oauth

import (
	"Blog-API/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	googleAPI "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type oauthService struct {
	googleConfig *oauth2.Config
	githubConfig *oauth2.Config
}

func NewOAuthService(googleCfg, githubCfg *oauth2.Config) domain.OAuthService {
	return &oauthService{
		googleConfig: googleCfg,
		githubConfig: githubCfg,
	}
}

func (s *oauthService) GetAuthURL(provider, state string) (string, error) {
	switch provider {
	case "google":
		return s.googleConfig.AuthCodeURL(state), nil
	case "github":
		return s.githubConfig.AuthCodeURL(state), nil
	default:
		return "", errors.New("unsupported oauth provider")
	}
}

func (s *oauthService) ExchangeCodeForToken(provider, code string) (*oauth2.Token, error) {
	switch provider {
	case "google":
		return s.googleConfig.Exchange(context.Background(), code)
	case "github":
		return s.githubConfig.Exchange(context.Background(), code)
	default:
		return nil, errors.New("unsupported oauth provider")
	}
}

func (s *oauthService) GetUserInfo(provider string, token *oauth2.Token) (oauthID, email, username string, err error) {
	switch provider {
	case "google":
		return s.getGoogleUserInfo(token)
	case "github":
		return s.getGitHubUserInfo(token)
	default:
		return "", "", "", errors.New("unsupported oauth provider")
	}
}

func (s *oauthService) getGoogleUserInfo(token *oauth2.Token) (string, string, string, error) {
	service, err := googleAPI.NewService(context.Background(), option.WithTokenSource(s.googleConfig.TokenSource(context.Background(), token)))
	if err != nil {
		return "", "", "", err
	}
	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return "", "", "", err
	}
	return userInfo.Id, userInfo.Email, userInfo.Name, nil
}

func (s *oauthService) getGitHubUserInfo(token *oauth2.Token) (string, string, string, error) {
	client := s.githubConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	var userInfo struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", "", "", err
	}
	if userInfo.Email == "" {
		return "", "", "", errors.New("github email is private or not set")
	}
	return fmt.Sprint(userInfo.ID), userInfo.Email, userInfo.Login, nil
}
