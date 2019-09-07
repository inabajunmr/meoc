package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/inabajunmr/meoc/oauth2"
)

type HttpRequest struct {
	Method string
	URI    string
}

func Call(httpRequest HttpRequest, oauth2Profile string) {

	// Get Access Token
	// TODO authentication info from file

	token := oauth2.GetAccessToken(oauth2Profile)
	fmt.Println(token)

	// Set Access Token for request
	client := new(http.Client)
	req, _ := http.NewRequest(httpRequest.Method, httpRequest.URI, nil)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Accept", "application/json")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}
