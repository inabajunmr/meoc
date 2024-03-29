package oauth2

import (
	"errors"
	"fmt"
	"os/user"
	"strings"

	"gopkg.in/ini.v1"
)

type OAuth2Config struct {
	ClientID                       string
	ClientSecret                   string
	TokenEndpoint                  string
	AuthorizationEndpoint          string
	GrantType                      string
	Scope                          string
	RedirectURI                    string
	Username                       string
	Password                       string
	TokenRequestParameters         map[string]string
	AuthorizationRequestParameters map[string]string
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefleshToken string `json:"reflesh_token"`
	Scope        string `json:"scope"`
}

func GetAccessToken(profile string) (*AccessToken, error) {

	usr, _ := user.Current()
	ini, _ := ini.Load(usr.HomeDir + "/.meoc/config")

	// TODO get AccessToken from cache
	config := OAuth2Config{
		ClientID:                       ini.Section(profile).Key("client_id").Value(),
		ClientSecret:                   ini.Section(profile).Key("client_secret").Value(),
		AuthorizationEndpoint:          ini.Section(profile).Key("authorization_endpoint").Value(),
		TokenEndpoint:                  ini.Section(profile).Key("token_endpoint").Value(),
		GrantType:                      ini.Section(profile).Key("grant_type").Value(),
		RedirectURI:                    ini.Section(profile).Key("redirect_uri").Value(),
		Scope:                          ini.Section(profile).Key("scope").Value(),
		Username:                       ini.Section(profile).Key("username").Value(),
		Password:                       ini.Section(profile).Key("password").Value(),
		TokenRequestParameters:         map[string]string{},
		AuthorizationRequestParameters: map[string]string{}}

	// add token request additional parameter
	for _, key := range ini.Section(profile).Keys() {
		if strings.HasPrefix(key.Name(), "token_request_p_") {
			config.TokenRequestParameters[key.Name()[16:]] = key.Value()
		}

		if strings.HasPrefix(key.Name(), "authorization_request_p_") {
			config.TokenRequestParameters[key.Name()[24:]] = key.Value()
		}
	}

	switch config.GrantType {
	case "client_credentials":
		return getTokenByClientCredentials(config)
	case "authorization_code":
		return getTokenByAuthorizationCode(config)
	case "password":
		return getTokenByResourceOwnersPasswordCredentials(config)
	case "implicit":
		return getTokenByImplicit(config)
	default:
		return nil, errors.New(fmt.Sprintf("Grant Type:%s is not supported.", config.GrantType))
	}
}
