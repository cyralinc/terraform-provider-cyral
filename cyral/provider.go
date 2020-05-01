package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
type auth0TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

// Provider defines and initializes the Cyral provider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth0_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auth0_client_id": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth0_client_secret": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth0_audience": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"control_plane": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cyral_repository": resourceCyralRepository(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := &Config{
		auth0Domain:       d.Get("auth0_domain").(string),
		auth0ClientID:     d.Get("auth0_client_id").(string),
		auth0ClientSecret: d.Get("auth0_client_secret").(string),
		auth0Audience:     d.Get("auth0_audience").(string),
		controlPlane:      d.Get("control_plane").(string),
		terraformVersion:  terraformVersion,
	}

	token, err := readTokenInfo(config.auth0Domain, config.auth0ClientID,
		config.auth0ClientSecret, config.auth0Audience)
	if err != nil {
		return nil, nil
	}

	config.token = token.AccessToken
	config.tokenType = token.TokenType

	return config, nil
}

func readTokenInfo(domain, clientID, clientSecret, audience string) (auth0TokenResponse, error) {
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
