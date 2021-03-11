package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Auth0TokenResponse represents the payload with the token response from Auth0.
type Auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// Auth0TokenRequest represents the payload for token requests to Auth0.
type Auth0TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

// Client stores data for all existing resources. Also, this is
// the struct that is passed along resources CRUD operations.
type Client struct {
	Token        string
	TokenType    string
	ControlPlane string
}

// NewClient configures and returns a fully initialized Client.
func NewClient(clientID, clientSecret, auth0Domain, auth0Audience,
	controlPlane string, keycloakProvider bool) (*Client, error) {

	if !keycloakProvider {
		token, err := getAuth0Token(auth0Domain, clientID, clientSecret, auth0Audience)
		if err != nil {
			return nil, err
		}

		return &Client{
			ControlPlane: controlPlane,
			Token:        token.AccessToken,
			TokenType:    token.TokenType,
		}, nil
	}
	return nil, fmt.Errorf("unsupported auth provider: keycloak. Please set 'auth_provider = \"auth0\"")
}

func getAuth0Token(domain, clientID, clientSecret, audience string) (Auth0TokenResponse, error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)
	audienceURL := fmt.Sprintf("https://%s", audience)
	tokenReq := Auth0TokenRequest{
		Audience:     audienceURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	payloadBytes, err := json.Marshal(tokenReq)
	if err != nil {
		return Auth0TokenResponse{}, fmt.Errorf("failed to encode readToken payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return Auth0TokenResponse{}, fmt.Errorf("unable to create auth0 request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Auth0TokenResponse{}, fmt.Errorf("unable execute auth0 request; err: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Auth0 requisition fail. Response status code %d.\n", res.StatusCode)
		return Auth0TokenResponse{}, fmt.Errorf(msg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Auth0TokenResponse{}, fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	token := Auth0TokenResponse{}
	if err := json.Unmarshal(body, &token); err != nil {
		return Auth0TokenResponse{}, fmt.Errorf("unable to get access token from json; err: %v", err)
	}

	return token, nil
}

// DoRequest calls the httpMethod informed and delivers the resourceData as a payload,
// filling the response parameter (if not nil) with the response body.
func (c *Client) DoRequest(url, httpMethod string, resourceData interface{}) ([]byte, error) {
	log.Printf("[DEBUG] Resource info: %#v", resourceData)
	log.Printf("[DEBUG] %s URL: %s", httpMethod, url)

	var req *http.Request

	if resourceData != nil {
		payloadBytes, err := json.Marshal(resourceData)
		if err != nil {
			return nil, fmt.Errorf("failed to encode payload: %v", err)
		}
		payload := string(payloadBytes)
		log.Printf("[DEBUG] %s payload: %s", httpMethod, payload)
		if req, err = http.NewRequest(httpMethod, url, strings.NewReader(payload)); err != nil {
			return nil, fmt.Errorf("unable to create request; err: %v", err)
		}
	} else {
		var err error
		if req, err = http.NewRequest(httpMethod, url, nil); err != nil {
			return nil, fmt.Errorf("unable to create request; err: %v", err)
		}
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))

	log.Printf("[DEBUG] Executing %s", httpMethod)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request. Check the control plane address; err: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusConflict ||
		strings.Contains(strings.ToLower(res.Status), "already exists") {
		return nil, fmt.Errorf("resource possibly exists in the control plane. Response status: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read data from request body; err: %v", err)
	}
	log.Printf("[DEBUG] Response body: %s", string(body))

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("resource not found; %v", resourceData)
	} else if res.StatusCode == http.StatusConflict {
		return nil, fmt.Errorf("resource conflict; status code: %d; body: %q",
			res.StatusCode, res.Body)
	} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error executing request; status code: %d; body: %q",
			res.StatusCode, res.Body)
	}

	return body, nil
}
