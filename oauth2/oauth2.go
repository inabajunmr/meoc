package oauth2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	TokenRequestParameters         map[string]string
	AuthorizationRequestParameters map[string]string
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	RefleshToken string `json:"reflesh_token"`
	Scope        string `json:"scope"`
}

func GetAccessToken(profile string) AccessToken {

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
	default:
		return AccessToken{} // TODO exception
	}
}

func getTokenByClientCredentials(config OAuth2Config) AccessToken {
	form := url.Values{}
	form.Add("client_id", config.ClientID)
	form.Add("client_secret", config.ClientSecret)
	form.Add("grant_type", config.GrantType)
	form.Add("scope", config.Scope)

	for key, value := range config.TokenRequestParameters {
		form.Add(key, value)
	}

	body := strings.NewReader(form.Encode())
	client := new(http.Client)
	req, _ := http.NewRequest("POST", config.TokenEndpoint, body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))

	// token request
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	tokenResponse := AccessToken{}
	json.Unmarshal(byteArray, &tokenResponse) // TODO error

	return tokenResponse
}

func getTokenByAuthorizationCode(config OAuth2Config) AccessToken {
	// compile authorization request uri
	authReq, _ := http.NewRequest("GET", config.AuthorizationEndpoint, nil)
	authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	query := authReq.URL.Query()
	query.Add("client_id", config.ClientID)
	query.Add("client_secret", config.ClientSecret)
	query.Add("response_type", "code")
	query.Add("scope", config.Scope)
	query.Add("redirect_uri", config.RedirectURI)
	for key, value := range config.AuthorizationRequestParameters {
		query.Add(key, value)
	}

	authReq.URL.RawQuery = query.Encode()

	// user access uri by browser
	fmt.Println("Access to ", authReq.URL.String())
	// user input auth code by terminal
	fmt.Println("Authentication code:")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	code := stdin.Text()

	// token request
	form := url.Values{}
	form.Add("client_id", config.ClientID)
	form.Add("client_secret", config.ClientSecret)
	form.Add("grant_type", config.GrantType)
	form.Add("code", code)
	form.Add("redirect_uri", config.RedirectURI)
	form.Add("response_type", "code")
	for key, value := range config.TokenRequestParameters {
		form.Add(key, value)
	}

	body := strings.NewReader(form.Encode())
	client := new(http.Client)
	tokenReq, _ := http.NewRequest("POST", config.TokenEndpoint, body)
	tokenReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.Header.Add("Accept", "application/json")

	dump, _ := httputil.DumpRequestOut(tokenReq, true)
	fmt.Println(string(dump))

	resp, err := client.Do(tokenReq)
	fmt.Println(err)

	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	tokenResponse := AccessToken{}
	json.Unmarshal(byteArray, &tokenResponse) // TODO error

	return tokenResponse
}
