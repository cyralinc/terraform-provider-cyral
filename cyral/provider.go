package cyral

import (
	"context"
	"fmt"
	"os"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	var diags diag.Diagnostics

	auth0ClientID, err := getEnv("AUTH0_CLIENT_ID")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read environment variable",
			Detail:   err.Error(),
		})

		return nil, diags
	}
	auth0ClientSecret, err := getEnv("AUTH0_CLIENT_SECRET")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read environment variable",
			Detail:   err.Error(),
		})

		return nil, diags
	}
	auth0Domain := d.Get("auth0_domain").(string)
	auth0Audience := d.Get("auth0_audience").(string)
	controlPlane := d.Get("control_plane").(string)

	c, err := client.NewClient(auth0ClientID, auth0ClientSecret, auth0Domain, auth0Audience,
		controlPlane)
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

func getEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("missing environment variable: %s", key)
}
