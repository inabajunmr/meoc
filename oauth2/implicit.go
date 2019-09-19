package oauth2

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func getTokenByImplicit(config OAuth2Config) (*AccessToken, error) {
	// compile authorization request uri
	authReq, _ := http.NewRequest("GET", config.AuthorizationEndpoint, nil)
	authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	query := authReq.URL.Query()
	query.Add("scope", config.Scope)
	query.Add("client_id", config.ClientID)
	query.Add("client_secret", config.ClientSecret)
	query.Add("response_type", "token")
	query.Add("redirect_uri", config.RedirectURI)
	for key, value := range config.AuthorizationRequestParameters {
		query.Add(key, value)
	}

	authReq.URL.RawQuery = query.Encode()

	// user access uri by browser
	fmt.Println("Access to ", authReq.URL.String())

	// user input hash by terminal
	fmt.Println("Hash:")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	hash := stdin.Text()

	// parse hash for constructing token
	p, err := url.Parse("http://example.com?" + hash)
	if err != nil {
		return nil, err
	}
	hashQuery := p.Query()
	var expiresIn int
	if hashQuery["expires_in"] != nil {
		expiresIn, _ = strconv.Atoi(hashQuery["expires_in"][0])
	}

	var scope string
	if hashQuery["scope"] != nil {
		scope = hashQuery["scope"][0]
	}

	tokenResponse := AccessToken{AccessToken: hashQuery["access_token"][0],
		TokenType: hashQuery["token_type"][0],
		ExpiresIn: expiresIn,
		Scope:     scope}

	return &tokenResponse, nil
}
