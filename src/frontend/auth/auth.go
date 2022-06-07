package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"frontend/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	REQUEST_STATE  = "velocimodel_state"
	CODE_CHALLENGE = "velocimodel-sha256-code-challenge"
)

var (
	OauthConfig  oauth2.Config
	ClientConfig clientcredentials.Config
)

func LoadOauthConfig() {
	OauthConfig = oauth2.Config{
		ClientID:     config.Config.Oauth.ClientID,
		ClientSecret: config.Config.Oauth.ClientSecret,
		Scopes:       []string{"all"},
		RedirectURL:  fmt.Sprintf("http://%s:9000/auth/redirect", config.Config.ExternalHTTPHost),
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Config.Oauth.AuthServerExternalURL + "/oauth/authorize",
			TokenURL: config.Config.Oauth.AuthServerInternalURL + "/oauth/token",
		},
	}
	ClientConfig = clientcredentials.Config{
		ClientID:     config.Config.Oauth.ClientID,
		ClientSecret: config.Config.Oauth.ClientSecret,
		TokenURL:     config.Config.Oauth.AuthServerInternalURL + "/oauth/token",
	}
}

func GetToken() (*oauth2.Token, error) {
	token, err := ClientConfig.Token(context.Background())
	if err != nil {
		return nil, err
	}
	return token, nil
}

func HandleRedirect(c *gin.Context) {
	state := c.Query("state")
	if state != REQUEST_STATE {
		c.Status(http.StatusBadRequest)
		c.Error(errors.New("state invlid"))
		return
	}
	code := c.Query("code")
	if code == "" {
		c.Status(http.StatusBadRequest)
		c.Error(errors.New("code not found"))
		return
	}
	token, err := OauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", CODE_CHALLENGE))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err)
		return
	}

	c.SetCookie("access_token", token.AccessToken, 3600, "/", config.Config.ExternalHTTPHost, true, false)
	c.SetCookie("refresh_token", token.RefreshToken, 3600, "/", config.Config.ExternalHTTPHost, true, false)

	c.Redirect(http.StatusFound, "/ui/dashboard")
}

func HandleLogin(c *gin.Context) {
	url := OauthConfig.AuthCodeURL(REQUEST_STATE,
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(CODE_CHALLENGE)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	c.Redirect(http.StatusFound, url)
}

func HandleLogout(c *gin.Context) {
	c.SetCookie("access_token", "", 0, "/", config.Config.ExternalHTTPHost, true, false)
	c.SetCookie("refresh_token", "", 0, "/", config.Config.ExternalHTTPHost, true, false)
	c.Redirect(http.StatusFound, "/auth/login")
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func generateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}
