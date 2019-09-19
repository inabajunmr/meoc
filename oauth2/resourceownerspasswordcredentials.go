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

func getTokenByResourceOwnersPasswordCredentials(config OAuth2Config) (*AccessToken, error) {
	// compile authorization request uri
	authReq, err := http.NewRequest("GET", config.AuthorizationEndpoint, nil)
	if err != nil {
		return nil, err
	}
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
	tokenReq, err := http.NewRequest("POST", config.TokenEndpoint, body)
	if err != nil {
		return nil, err
	}
	tokenReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.Header.Add("Accept", "application/json")

	dump, err := httputil.DumpRequestOut(tokenReq, true)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dump))

	resp, err := client.Do(tokenReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokenResponse := AccessToken{}
	if err = json.Unmarshal(byteArray, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
