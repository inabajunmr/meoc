package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type OAuth2Config struct {
	ClientId      string
	ClientSecret  string
	TokenEndpoint string
	GrantType     string
	Scope         string
}

func GetAccessToken(config OAuth2Config) string {
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
