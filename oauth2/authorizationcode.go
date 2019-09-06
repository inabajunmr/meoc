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
	"strings"
)

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
