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
	ClientId      string
	ClientSecret  string
	TokenEndpoint string
	GrantType     string
	Scope         string
}

func GetAccessToken(profile string) string {
	usr, _ := user.Current()
	ini, _ := ini.Load(usr.HomeDir + "/.meoc/config")

	config := OAuth2Config{
		ClientId:      ini.Section(profile).Key("client_id").String(),
		ClientSecret:  ini.Section(profile).Key("client_secret").String(),
		TokenEndpoint: ini.Section(profile).Key("token_endpoint").String(),
		GrantType:     ini.Section(profile).Key("grant_type").String(),
		Scope:         ini.Section(profile).Key("scope").String()}

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
