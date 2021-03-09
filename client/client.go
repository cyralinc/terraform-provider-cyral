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

// CreateResource calls a POST api to create a given resource and returns
// the unmarshalled body based on the response interface provided.
func (c *Client) CreateResource(resourceData interface{}, route string, response interface{}) error {
	log.Printf("[DEBUG] Resource info: %#v", resourceData)

	url := fmt.Sprintf("https://%s/v1/%s", c.ControlPlane, route)
	log.Printf("[DEBUG] POST URL: %s", url)
	payloadBytes, err := json.Marshal(resourceData)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %v", err)
	}

	log.Printf("[DEBUG] POST payload: %s", string(payloadBytes))
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))

	log.Printf("[DEBUG] Executing POST")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to execute request. Check the control plane address; err: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusConflict ||
		strings.Contains(strings.ToLower(res.Status), "already exists") {
		return fmt.Errorf("resource already exists in control plane")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read data from request body; err: %v", err)
	}
	log.Printf("[DEBUG] Response body: %s", string(body))

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response from CreateResource request; status code: %d; body: %q",
			res.StatusCode, body)
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("unable to unmarshall json; err: %v", err)
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	return nil
}

// ReadResource calls a GET api to read a given url and returns
// the unmarshalled body based on the response interface provided.
func (c *Client) ReadResource(resourceID, url string, response interface{}) error {
	log.Printf("[DEBUG] GET URL: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create new request; err: %v", err)
	}

	log.Printf("[DEBUG] Executing GET")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	log.Printf("[DEBUG] GET request: %#v", req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to execute request at ReadResource."+
			" Check the control plane address; err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read data from request body at ReadResource; err: %v", err)
	}
	log.Printf("[DEBUG] Response body: %s", string(body))

	// Not an error, nor any data was found
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found; id: %s", resourceID)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from ReadResource; status code: %d; body: %q", res.StatusCode, res.Body)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("unable to unmarshall response, err: %v", err)
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	return nil
}

// UpdateResource runs a PUT request on a given url and resource data.
func (c *Client) UpdateResource(resourceData interface{}, url string) error {
	payloadBytes, err := json.Marshal(resourceData)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to execute request. Check the control plane address; err: %v", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found; %v", resourceData)
	} else if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("resource conflicts in control plane; status code: %d; body: %q",
			res.StatusCode, res.Body)
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response; status code: %d; body: %q",
			res.StatusCode, res.Body)
	}

	return nil
}

// DeleteResource runs a DELETE request on a given url.
func (c *Client) DeleteResource(url string) error {
	log.Printf("[DEBUG] DELETE URL: %s", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create 'delete repo' request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable execute 'delete repo' request. Check the control plane address; err: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from 'delete repo' request; status code: %d; body: %q", res.StatusCode, res.Body)
	}

	return nil
}
