package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	u "net/url"
	"strings"
)

// Auth0TokenRequest represents the payload for token requests to Auth0.
type Auth0TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

// TokenResponse represents the payload with the token response from Auth0.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
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
	log.Printf("[DEBUG] Init NewClient")

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
	token, err := getKeycloakToken(controlPlane, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] token.TokenType: %s", token.TokenType)
	log.Printf("[DEBUG] token.AccessToken: %s", token.AccessToken)

	log.Printf("[DEBUG] End NewClient")

	return &Client{
		ControlPlane: controlPlane,
		Token:        token.AccessToken,
		TokenType:    token.TokenType,
	}, nil
}

func getAuth0Token(domain, clientID, clientSecret, audience string) (TokenResponse, error) {
	log.Printf("[DEBUG] Init getAuth0Token")

	url := fmt.Sprintf("https://%s/oauth/token", domain)
	audienceURL := fmt.Sprintf("https://%s", audience)
	tokenReq := Auth0TokenRequest{
		Audience:     audienceURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	log.Printf("[DEBUG] url: %s", url)
	log.Printf("[DEBUG] payload: %v", tokenReq)

	payloadBytes, err := json.Marshal(tokenReq)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to encode readToken payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return TokenResponse{}, fmt.Errorf("unable to create auth0 request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("unable execute auth0 request; err: %v", err)
	}
	defer res.Body.Close()
	log.Printf("[DEBUG] body: %v", res.Body)
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Auth0 requisition fail. Response status code %d. Response body: %v",
			res.StatusCode, res.Body)
		return TokenResponse{}, fmt.Errorf(msg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	token := TokenResponse{}
	if err := json.Unmarshal(body, &token); err != nil {
		return TokenResponse{}, fmt.Errorf("unable to get access token from json; err: %v", err)
	}

	log.Printf("[DEBUG] End getAuth0Token")

	return token, nil
}

func getKeycloakToken(controlPlane, clientID, clientSecret string) (TokenResponse, error) {
	log.Printf("[DEBUG] Init getKeycloakToken")
	url := fmt.Sprintf("https://%s/v1/users/oidc/token", controlPlane)

	log.Printf("[DEBUG] url: %s", url)
	log.Printf("[DEBUG] clientId: %s ; clientSecret: %s", clientID, clientSecret)

	data := u.Values{}
	data.Set("clientId", clientID)
	data.Set("clientSecret", clientSecret)
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return TokenResponse{}, fmt.Errorf("unable to create keycloak request; err: %v", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("unable execute keycloak request; err: %v", err)
	}
	defer res.Body.Close()
	log.Printf("[DEBUG] body: %v", res.Body)
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("keycloak requisition fail. Response status code %d. Response body: %v",
			res.StatusCode, res.Body)
		return TokenResponse{}, fmt.Errorf(msg)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	token := TokenResponse{}
	if err := json.Unmarshal(body, &token); err != nil {
		return TokenResponse{}, fmt.Errorf("unable to get access token from json; err: %v", err)
	}

	log.Printf("[DEBUG] End getKeycloakToken")

	return token, nil
}

// DoRequest calls the httpMethod informed and delivers the resourceData as a payload,
// filling the response parameter (if not nil) with the response body.
func (c *Client) DoRequest(url, httpMethod string, resourceData interface{}) ([]byte, error) {
	log.Printf("[DEBUG] Init DoRequest")
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
	// The TokenType returned by getKeycloakToken is "bearer", but if we use it here we
	// will get error "Failed to get roles: tokenstring should not contain 'bearer '\n".
	// If we change it to "Bearer" it works normally.
	// See: https://cyralinc.atlassian.net/browse/ENG-4408
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	//req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))

	log.Printf("[DEBUG] Executing %s", httpMethod)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request. Check the control plane address; err: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusConflict ||
		(httpMethod == http.MethodPost && strings.Contains(strings.ToLower(res.Status), "already exists")) {
		return nil, fmt.Errorf("resource possibly exists in the control plane. Response status: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read data from request body; err: %v", err)
	}
	log.Printf("[DEBUG] Request: %#v", req)
	log.Printf("[DEBUG] Response body: %s", string(body))

	var er error
	if res.StatusCode == http.StatusNotFound {
		er = fmt.Errorf("resource not found; %v", resourceData)
	} else if res.StatusCode == http.StatusConflict {
		er = fmt.Errorf("resource conflict; status code: %d; body: %q",
			res.StatusCode, body)
	} else if res.StatusCode != http.StatusOK &&
		(httpMethod == http.MethodPost && res.StatusCode != http.StatusCreated) {
		er = fmt.Errorf("error executing request; status code: %d; body: %q",
			res.StatusCode, body)
	} else if httpMethod == http.MethodDelete && res.StatusCode == 500 &&
		!strings.Contains(strings.ToLower(string(body)), "does not exist") {
		// Ignore http 500 error codes that informs that the resource was not found
		// in the server. In these cases, will consider as deleted. In such cases,
		// the correct behavior would be that the API would return http 410 (Gone).
		er = fmt.Errorf("error executing request; status code: %d; body: %q",
			res.StatusCode, body)
	}

	log.Printf("[DEBUG] End DoRequest")

	return body, er
}
