package client

import (
	"fmt"

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
	// Set TODO Access Token for request

	// req, _ := http.NewRequest(httpRequest.Method, httpRequest.URI, nil)

	// resp, _ := client.Do(req)
	// defer resp.Body.Close()

	// byteArray, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(byteArray))
}
