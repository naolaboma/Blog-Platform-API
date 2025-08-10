package controllers

import (
	"Blog-API/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	userUseCase  domain.UserUseCase
	oauthService domain.OAuthService
	stateSecret  string
}

func NewOAuthHandler(userUseCase domain.UserUseCase, oauthService domain.OAuthService, stateSecret string) *OAuthHandler {
	return &OAuthHandler{
		userUseCase:  userUseCase,
		oauthService: oauthService,
		stateSecret:  stateSecret,
	}
}

func (h *OAuthHandler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")
	//generate random state for csfr protection
	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)
	//set state in a secure httponly cookie
	c.SetCookie("oauthstate", state, int(10*time.Minute.Seconds()), "/", "localhost", true, true)

	// get the redirect url from the service
	url, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	//redirect the user's browser to the provider's login page
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func (h *OAuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	// get the state from the cookie and query Param
	storedState, err := c.Cookie("oauthstate")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "missing oaut state cookie"})
		return
	}
	queryState := c.Query("state")
	// clear the cookie immidiately
	c.SetCookie("oauthstate", "", -1, "/", "localhost", true, true)

	//get the authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "oauth code is missing"})
		return
	}

	loginResponse, err := h.userUseCase.OAuthLogin(provider, queryState, code, storedState)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}
