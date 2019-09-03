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
	ClientId              string
	ClientSecret          string
	TokenEndpoint         string
	AuthorizationEndpoint string
	GrantType             string
	Scope                 string
	RedirectUri           string
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
		ClientId:              ini.Section(profile).Key("client_id").String(),
		ClientSecret:          ini.Section(profile).Key("client_secret").String(),
		AuthorizationEndpoint: ini.Section(profile).Key("authorization_endpoint").String(),
		TokenEndpoint:         ini.Section(profile).Key("token_endpoint").String(),
		GrantType:             ini.Section(profile).Key("grant_type").String(),
		RedirectUri:           ini.Section(profile).Key("redirect_uri").String(),
		Scope:                 ini.Section(profile).Key("scope").String()}

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
	tokenResponse := AccessToken{}
	json.Unmarshal(byteArray, &tokenResponse) // TODO error
	return tokenResponse
}

func getTokenByAuthorizationCode(config OAuth2Config) AccessToken {
	// compile authorization request uri
	authReq, _ := http.NewRequest("GET", config.AuthorizationEndpoint, nil)
	authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	query := authReq.URL.Query()
	query.Add("client_id", config.ClientId)
	query.Add("client_secret", config.ClientSecret)
	query.Add("response_type", "code")
	query.Add("scope", config.Scope)
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
	form.Add("client_id", config.ClientId)
	form.Add("client_secret", config.ClientSecret)
	form.Add("grant_type", config.GrantType)
	form.Add("redirect_uri", config.RedirectUri)
	form.Add("code", code)

	body := strings.NewReader(form.Encode())
	client := new(http.Client)
	tokenReq, _ := http.NewRequest("POST", config.TokenEndpoint, body)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.Header.Set("Accept", "application/json")

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
