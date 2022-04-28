package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	authServerInternalURL = "http://auth-manager:9005"
	authServerExternalURL = "http://localhost:9005"
	REQUEST_STATE         = "velocimodel_state"
	CODE_CHALLENGE        = "velocimodel-sha256-code-challenge"
)

var (
	conf = oauth2.Config{
		ClientID:     "velocimodel-frontend",
		ClientSecret: "abc123DEFghijklmnop4567rstuvwxyzZYXWUT8910SRQPOnmlijhoauthplaygroundapplication",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9000/auth/redirect",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerExternalURL + "/oauth/authorize",
			TokenURL: authServerInternalURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func HandleRedirect(c *gin.Context) {
	state := c.Query("state")
	if state != REQUEST_STATE {
		c.Status(http.StatusBadRequest)
		c.Error(errors.New("State invlid"))
		return
	}
	code := c.Query("code")
	if code == "" {
		c.Status(http.StatusBadRequest)
		c.Error(errors.New("Code not found"))
		return
	}
	token, err := conf.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", CODE_CHALLENGE))
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err)
		return
	}

	c.SetCookie("access_token", token.AccessToken, 3600, "/", "localhost", true, false)
	c.SetCookie("refresh_token", token.RefreshToken, 3600, "/", "localhost", true, false)

	c.Redirect(http.StatusFound, "/ui/dashboard")
}

func HandleLogin(c *gin.Context) {
	url := conf.AuthCodeURL(REQUEST_STATE,
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(CODE_CHALLENGE)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	log.Printf("REDIRECT URL LOGIN: %v", url)
	c.Redirect(http.StatusFound, url)
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
