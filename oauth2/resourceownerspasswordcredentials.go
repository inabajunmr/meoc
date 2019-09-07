package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func getTokenByResourceOwnersPasswordCredentials(config OAuth2Config) AccessToken {
	// compile authorization request uri
	authReq, _ := http.NewRequest("GET", config.AuthorizationEndpoint, nil)
	authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	form := url.Values{}
	form.Add("scope", config.Scope)
	form.Add("grant_type", "password")
	form.Add("password", config.Password)
	form.Add("username", config.Username)
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
