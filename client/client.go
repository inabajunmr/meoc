package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/inabajunmr/meoc/oauth2"
)

type HttpRequest struct {
	Method  string
	URI     string
	Headers []Header
}

type Header struct {
	Name  string
	Value string
}

func Call(httpRequest HttpRequest, oauth2Profile string) error {
	// Get Access Token
	token, err := oauth2.GetAccessToken(oauth2Profile)
	if err != nil {
		return err
	}

	// Set Access Token for request
	client := new(http.Client)
	req, err := http.NewRequest(httpRequest.Method, httpRequest.URI, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Accept", "application/json")
	for _, v := range httpRequest.Headers {
		req.Header.Add(v.Name, v.Value)
	}

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	fmt.Println(string(dump))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(byteArray))
	return nil
}
