package cyral

import (
	"context"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keycloak = "keycloak"
	auth0    = "auth0"
)

// Provider defines and initializes the Cyral provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_provider": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  keycloak,
			},
			"auth0_audience": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auth0_domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CYRAL_TF_CLIENT_ID", nil),
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CYRAL_TF_CLIENT_SECRET", nil),
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

	keycloakProvider := d.Get("auth_provider").(string) == keycloak

	clientID := d.Get("client_id").(string)
	if clientID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read environment variable",
			Detail:   "Unable to read environment variable CYRAL_TF_CLIENT_ID",
		})

		return nil, diags
	}
	clientSecret := d.Get("client_secret").(string)
	if clientSecret == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read environment variable",
			Detail:   "Unable to read environment variable CYRAL_TF_CLIENT_SECRET",
		})

		return nil, diags
	}
	auth0Domain := d.Get("auth0_domain").(string)
	auth0Audience := d.Get("auth0_audience").(string)
	controlPlane := d.Get("control_plane").(string)

	c, err := client.NewClient(clientID, clientSecret, auth0Domain, auth0Audience,
		controlPlane, keycloakProvider)
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
