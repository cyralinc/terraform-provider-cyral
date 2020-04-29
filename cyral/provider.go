package cyral

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider defines and initializes the Cyral provider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		/*Schema: map[string]*schema.Schema{
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
		},*/
		ResourcesMap: map[string]*schema.Resource{
			"cyral_repository": resourceCyralRepository(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cyral_data_source_auth0": dataSourceCyralAuth0(),
		},
	}

	/*provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}*/

	return provider
}

/*
func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	domain := d.Get("auth0_domain").(string)
	clientID := d.Get("auth0_client_id").(string)
	clientSecret := d.Get("auth0_client_secret").(string)
	jwtToken, _ := readJWTToken(domain, clientID, clientSecret)

	config := Config{
		Auth0Domain:       domain,
		Auth0ClientID:     clientID,
		Auth0ClientSecret: clientSecret,
		JWTToken:          jwtToken,
		terraformVersion:  terraformVersion,
	}

	return config, nil
}

func readJWTToken(domain, clientID, clientSecret string) (string, error) {

	url := fmt.Sprintf("https://%s/oauth/token", domain)
	audienceURL := fmt.Sprintf("https://%s/api", domain)

	payloadStr := fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"client_credentials\"}",
		clientID, clientSecret, audienceURL)

	payload := strings.NewReader(payloadStr)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {

	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {

	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Printf("[INFO] ========\nres")
	log.Printf(fmt.Sprintf("[INFO] %s", res))

	log.Printf("[INFO] ========\nbody")
	log.Printf(fmt.Sprintf("[INFO] %s", string(body)))

	return "", nil
}
*/
