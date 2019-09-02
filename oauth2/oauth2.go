package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/user"
	"strings"

	"gopkg.in/ini.v1"
)

type OAuth2Config struct {
	ClientId              string
	ClientSecret          string
	TokenEndpoint         string
	AuthorizationEndpoint string
	GrantType             string
	Scope                 string
}

func GetAccessToken(profile string) string {

	usr, _ := user.Current()
	ini, _ := ini.Load(usr.HomeDir + "/.meoc/config")

	// TODO get AccessToken from cache
	config := OAuth2Config{
		ClientId:              ini.Section(profile).Key("client_id").String(),
		ClientSecret:          ini.Section(profile).Key("client_secret").String(),
		AuthorizationEndpoint: ini.Section(profile).Key("authorization_endpoint").String(),
		TokenEndpoint:         ini.Section(profile).Key("token_endpoint").String(),
		GrantType:             ini.Section(profile).Key("grant_type").String(),
		Scope:                 ini.Section(profile).Key("scope").String()}

	switch config.GrantType {
	case "client_credentials":
		return getTokenByClientCredentials(config)
	case "authorization_code":
		return getTokenByAuthorizationCode(config)
	default:
		return "" // TODO exception
	}
}

func getTokenByClientCredentials(config OAuth2Config) string {
	form := url.Values{}
	form.Add("client_id", config.ClientId)
	form.Add("client_secret", config.ClientSecret)
	form.Add("grant_type", config.GrantType)
	form.Add("scope", config.Scope)
	form.Add("audience", "http://example.com")

	body := strings.NewReader(form.Encode())
	client := new(http.Client)
	req, _ := http.NewRequest("POST", config.TokenEndpoint, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	return string(byteArray)
}

func getTokenByAuthorizationCode(config OAuth2Config) string {
	req, _ := http.NewRequest("GET", config.AuthorizationEndpoint, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	query := req.URL.Query()
	query.Add("client_id", config.ClientId)
	query.Add("client_secret", config.ClientSecret)
	query.Add("response_type", "code")
	query.Add("scope", config.Scope)
	req.URL.RawQuery = query.Encode()
	fmt.Println(req.URL.String)

	// get auth code
	// get token
	return "string(byteArray)"
}
