package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Config contains the provider configuration parameters stored
// during provider initialization.
type Config struct {
	Auth0Domain   string
	Auth0Audience string
	controlPlane  string

	terraformVersion string
}

// CyralClient stores data for all existing resources. Also, this is
// the struct that is passed along resources CRUD operations.
type CyralClient struct {
	Token        string
	TokenType    string
	ControlPlane string

	Repository Repository
}

// Repository struct stores data for resource repository.
type Repository struct {
	Name string
}

// Client configures and returns a fully initialized CyralClient.
func (c *Config) Client() (interface{}, error) {
	auth0ClientID, err := c.getEnv("AUTH0_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	auth0ClientSecret, err := c.getEnv("AUTH0_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	token, err := c.readTokenInfo(c.Auth0Domain, auth0ClientID, auth0ClientSecret, c.Auth0Audience)
	if err != nil {
		return nil, err
	}

	return &CyralClient{
		ControlPlane: c.controlPlane,
		Token:        token.AccessToken,
		TokenType:    token.TokenType,
		Repository:   Repository{},
	}, nil
}

func (c *Config) readTokenInfo(domain, clientID, clientSecret, audience string) (auth0TokenResponse, error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)
	audienceURL := fmt.Sprintf("https://%s", audience)
	tokenReq := auth0TokenRequest{
		Audience:     audienceURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	payloadBytes, err := json.Marshal(tokenReq)
	if err != nil {
		return auth0TokenResponse{}, fmt.Errorf("failed to encode readToken payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return auth0TokenResponse{}, fmt.Errorf("unable to create auth0 request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return auth0TokenResponse{}, fmt.Errorf("unable execute auth0 request; err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return auth0TokenResponse{}, fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	token := auth0TokenResponse{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return auth0TokenResponse{}, fmt.Errorf("unable to get access token from json; err: %v", err)
	}

	return token, nil
}

func (c *Config) getEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("missing environment variable: %s", key)
}
