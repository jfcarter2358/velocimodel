package auth

import (
	"context"

	"service-manager/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	REQUEST_STATE  = "velocimodel_state"
	CODE_CHALLENGE = "velocimodel-sha256-code-challenge"
)

var (
	ClientConfig clientcredentials.Config
)

func LoadOauthConfig() {
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
