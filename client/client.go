package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	cc "golang.org/x/oauth2/clientcredentials"
)

// Client stores data for all existing resources. Also, this is
// the struct that is passed along resources CRUD operations.
type Client struct {
	ControlPlane string
	TokenSource  oauth2.TokenSource
	client       *http.Client
}

// NewClient configures and returns a fully initialized Client.
func NewClient(clientID, clientSecret, controlPlane string, tlsSkipVerify bool) (*Client, error) {
	log.Printf("[DEBUG] Init NewClient")

	if clientID == "" || clientSecret == "" || controlPlane == "" {
		return nil, fmt.Errorf("clientID, clientSecret and controlPlane must have non-empty values")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: tlsSkipVerify,
			},
		},
	}

	tokenConfig := cc.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("https://%s/v1/users/oidc/token", controlPlane),
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	tokenSource := tokenConfig.TokenSource(context.Background())

	log.Printf("[DEBUG] TokenSource: %v", tokenSource)
	log.Printf("[DEBUG] End NewClient")

	return &Client{
		ControlPlane: controlPlane,
		TokenSource:  tokenSource,
		client:       client,
	}, nil
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
	if c.TokenSource != nil {
		if token, err := c.TokenSource.Token(); err != nil {
			return nil, fmt.Errorf("unable to retrieve authorization token. error: %v", err)
		} else {
			log.Printf("[DEBUG] Token Type: %s", token.Type())
			log.Printf("[DEBUG] Access Token: %s", token.AccessToken)
			log.Printf("[DEBUG] Token Expiry: %s", token.Expiry)
			req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.Type(), token.AccessToken))
		}
	}

	log.Printf("[DEBUG] Executing %s", httpMethod)
	res, err := c.client.Do(req)
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
	log.Printf("[DEBUG] Response status code: %d", res.StatusCode)
	log.Printf("[DEBUG] Response body: %s", string(body))

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		err = fmt.Errorf("error executing %s request; status code: %d; body: %q",
			httpMethod, res.StatusCode, body)
	}

	log.Printf("[DEBUG] End DoRequest")

	return body, err
}
