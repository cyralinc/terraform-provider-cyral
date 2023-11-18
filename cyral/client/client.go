package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	cc "golang.org/x/oauth2/clientcredentials"
)

const redactedString = "**********"

const (
	EnvVarClientID      = "CYRAL_TF_CLIENT_ID"
	EnvVarClientSecret  = "CYRAL_TF_CLIENT_SECRET"
	EnvVarCPURL         = "CYRAL_TF_CONTROL_PLANE"
	EnvVarTLSSkipVerify = "CYRAL_TF_TLS_SKIP_VERIFY"
)

// Client stores data for all existing resources. Also, this is
// the struct that is passed along resources CRUD operations.
type Client struct {
	ControlPlane string
	TokenSource  oauth2.TokenSource
	client       *http.Client
}

// New configures and returns a fully initialized Client.
func New(clientID, clientSecret, controlPlane string, tlsSkipVerify bool) (*Client, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Init client.New")

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
	tokenSource := tokenConfig.TokenSource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("TokenSource: %v", tokenSource))
	tflog.Debug(ctx, "End client.New")

	return &Client{
		ControlPlane: controlPlane,
		TokenSource:  tokenSource,
		client:       client,
	}, nil
}

// DoRequest calls the httpMethod informed and delivers the resourceData as a payload,
// filling the response parameter (if not nil) with the response body.
func (c *Client) DoRequest(url, httpMethod string, resourceData interface{}) ([]byte, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Init DoRequest")
	tflog.Debug(ctx, fmt.Sprintf("Resource info: %#v", resourceData))
	tflog.Debug(ctx, fmt.Sprintf("%s URL: %s", httpMethod, url))
	var req *http.Request
	var err error
	if resourceData != nil {
		payloadBytes, err := json.Marshal(resourceData)
		if err != nil {
			return nil, fmt.Errorf("failed to encode payload: %v", err)
		}
		payload := string(payloadBytes)
		tflog.Debug(ctx, fmt.Sprintf("%s payload: %s", httpMethod, payload))
		if req, err = http.NewRequest(httpMethod, url, strings.NewReader(payload)); err != nil {
			return nil, fmt.Errorf("unable to create request; err: %v", err)
		}
	} else {
		if req, err = http.NewRequest(httpMethod, url, nil); err != nil {
			return nil, fmt.Errorf("unable to create request; err: %v", err)
		}
	}

	req.Header.Add("content-type", "application/json")
	token := &oauth2.Token{}
	if c.TokenSource != nil {
		if token, err = c.TokenSource.Token(); err != nil {
			return nil, fmt.Errorf("unable to retrieve authorization token. error: %v", err)
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Token Type: %s", token.Type()))
			tflog.Debug(ctx, fmt.Sprintf("Access Token: %s", redactContent(token.AccessToken)))
			tflog.Debug(ctx, fmt.Sprintf("Token Expiry: %s", token.Expiry))
			req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.Type(), token.AccessToken))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Executing %s", httpMethod))
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request. Check the control plane address; err: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusConflict ||
		(httpMethod == http.MethodPost && strings.Contains(strings.ToLower(res.Status), "already exists")) {
		return nil, NewHttpError(
			fmt.Sprintf("resource possibly exists in the control plane. Response status: %s", res.Status),
			res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, NewHttpError(
			fmt.Sprintf("unable to read data from request body; err: %v", err),
			res.StatusCode)
	}

	// Redact token before logging the request
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type(), redactContent(token.AccessToken)))

	tflog.Debug(ctx, fmt.Sprintf("Request: %#v", req))
	tflog.Debug(ctx, fmt.Sprintf("Response status code: %d", res.StatusCode))
	tflog.Debug(ctx, fmt.Sprintf("Response body: %s", string(body)))

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, NewHttpError(
			fmt.Sprintf("error executing %s request; status code: %d; body: %q",
				httpMethod, res.StatusCode, body),
			res.StatusCode)
	}

	tflog.Debug(ctx, "End DoRequest")

	return body, nil
}

func redactContent(content string) string {
	if content == "" {
		return content
	}
	return redactedString
}

func FromEnv() (*Client, error) {
	clientID, clientSecret, controlPlane, tlsSkipVerify, err :=
		getProviderConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("unable to create Cyral client: %w", err)
	}
	c, err := New(clientID, clientSecret, controlPlane,
		tlsSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("unable to create Cyral client: %w", err)
	}
	return c, nil
}

func getProviderConfigFromEnv() (
	clientID string,
	clientSecret string,
	controlPlane string,
	tlsSkipVerify bool,
	err error,
) {
	clientID = os.Getenv(EnvVarClientID)
	clientSecret = os.Getenv(EnvVarClientSecret)
	controlPlane = os.Getenv(EnvVarCPURL)
	tlsSkipVerifyStr := os.Getenv(EnvVarTLSSkipVerify)
	if tlsSkipVerifyStr != "" {
		tlsSkipVerify, err = strconv.ParseBool(tlsSkipVerifyStr)
		if err != nil {
			return "", "", "", false, fmt.Errorf("invalid value for "+
				"env var %q: %w", EnvVarTLSSkipVerify, err)
		}
	}
	return
}
