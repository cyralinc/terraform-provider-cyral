package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCyralAuth0() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCyralAuth0Read,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"client_secret": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceCyralAuth0Read(d *schema.ResourceData, m interface{}) error {
	domain := d.Get("domain").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	token, err := readJWTToken(domain, clientID, clientSecret)
	if err != nil {
		return err
	}

	d.Set("auth0_access_token", token)

	return nil
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
