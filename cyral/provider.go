package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
		Auth0Domain:       d.Get("auth0_domain").(string),
		Auth0ClientID:     d.Get("auth0_client_id").(string),
		Auth0ClientSecret: d.Get("auth0_client_secret").(string),
		terraformVersion:  terraformVersion,
	}

	token, err := readJWTToken(config.Auth0Domain, config.Auth0ClientID, config.Auth0ClientSecret)
	if err != nil {
		return nil, nil
	}

	config.JWTToken = token

	return config, nil
}

func readJWTToken(domain, clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://%s/oauth/token", domain)
	audienceURL := fmt.Sprintf("https://%s/api/v2/", domain)
	payloadStr := fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"client_credentials\"}",
		clientID, clientSecret, audienceURL)
	payload := strings.NewReader(payloadStr)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("unable to create auth0 request, err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable execute auth0 request, err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read data from request body, err: %v", err)
	}

	type Auth0Token struct {
		AccessToken string `json:"access_token"`
	}
	token := &Auth0Token{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return "", fmt.Errorf("unable to get access token from json, err: %v", err)
	}

	return token.AccessToken, nil
}
