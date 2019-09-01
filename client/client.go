package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type HttpRequest struct {
	Method string
	URI    string
}

func Call(httpRequest HttpRequest) {

	// Get Access Token
	grantType := "client_credentials"
	scope := ""
	tokenEndpoint := "https://meoc.auth0.com/oauth/token"
	clientSecret := "xx"
	clientId := "xx"
	form := url.Values{}
	form.Add("grant_type", grantType)
	form.Add("scope", scope)
	form.Add("client_id", clientId)
	form.Add("client_secret", clientSecret)
	form.Add("audience", "http://example.com")

	body := strings.NewReader(form.Encode())
	client := new(http.Client)
	req, _ := http.NewRequest("POST", tokenEndpoint, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))

	// Set TODO Access Token for request

	// req, _ := http.NewRequest(httpRequest.Method, httpRequest.URI, nil)

	// resp, _ := client.Do(req)
	// defer resp.Body.Close()

	// byteArray, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(byteArray))
}
