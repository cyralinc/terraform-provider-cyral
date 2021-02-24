package cyral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth0_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := &Config{
		Auth0Domain:   d.Get("auth0_domain").(string),
		Auth0Audience: d.Get("auth0_audience").(string),
		controlPlane:  d.Get("control_plane").(string),
	}

	var diags diag.Diagnostics

	var c, err = config.Client()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Cyral client",
			Detail:   err.Error(),
		})

		return nil, diags
	}

	return c, diags
}
