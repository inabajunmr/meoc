package client

import (
	"github.com/inabajunmr/meoc/oauth2"
)

type HttpRequest struct {
	Method string
	URI    string
}

func Call(httpRequest HttpRequest) {

	// Get Access Token
	// TODO authentication info from file
	config := oauth2.OAuth2Config{"xx",
		"xx",
		"https://meoc.auth0.com/oauth/token", "client_credentials", ""}
	oauth2.GetAccessToken(config)
	// Set TODO Access Token for request

	// req, _ := http.NewRequest(httpRequest.Method, httpRequest.URI, nil)

	// resp, _ := client.Do(req)
	// defer resp.Body.Close()

	// byteArray, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(byteArray))
}
