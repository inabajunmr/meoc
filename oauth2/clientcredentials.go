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

func getTokenByClientCredentials(config OAuth2Config) (*AccessToken, error) {
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
	req, err := http.NewRequest("POST", config.TokenEndpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dump))

	// token request
	resp, err := client.Do(req)
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
		fmt.Println(string(byteArray))
		return nil, err
	}

	return &tokenResponse, nil
}
