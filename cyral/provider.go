package cyral

import (
	"context"
	"fmt"

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
			"auth_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  keycloak,
			},
			"auth0_audience": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth0_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth0_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AUTH0_CLIENT_ID", nil),
				ConflictsWith: []string{"client_id"},
				Deprecated: "use provider variable 'client_id' or environment variable " +
					"'CYRAL_TF_CLIENT_ID' instead of 'auth0_client_id'",
			},
			"auth0_client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AUTH0_CLIENT_SECRET", nil),
				ConflictsWith: []string{"client_secret"},
				Deprecated: "use provider variable 'client_secret' or environment variable " +
					"'CYRAL_TF_CLIENT_SECRET' instead of 'auth0_client_secret'",
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"auth0_client_id"},
				DefaultFunc:   schema.EnvDefaultFunc("CYRAL_TF_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"auth0_client_secret"},
				DefaultFunc:   schema.EnvDefaultFunc("CYRAL_TF_CLIENT_SECRET", nil),
			},
			"control_plane": {
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
	keycloakProvider := d.Get("auth_provider").(string) == keycloak

	clientID, clientSecret, diags := getCredentials(d, keycloakProvider)
	if clientID == "" || clientSecret == "" {
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

func getCredentials(d *schema.ResourceData, keycloakProvider bool) (string, string, diag.Diagnostics) {
	var clientID, clientSecret string

	getVar := func(providerVar, envVar string, diags *diag.Diagnostics) string {
		value := d.Get(providerVar).(string)
		if value == "" {
			(*diags) = append((*diags), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read credentials",
				Detail:   fmt.Sprintf("use provider variable '%s' or environment variable '%s'", providerVar, envVar),
			})
		}
		return value
	}
	var diags diag.Diagnostics

	clientID = getVar("client_id", "CYRAL_TF_CLIENT_ID", &diags)
	clientSecret = getVar("client_secret", "CYRAL_TF_CLIENT_SECRET", &diags)

	// Backwards compatibility code to allow users to migrate to new variables and see
	// a deprecation warning. The code below must be removed in next versions.
	if !keycloakProvider && clientID == "" && clientSecret == "" {
		diags = nil
		clientID = getVar("auth0_client_id", "AUTH0_CLIENT_ID", &diags)
		clientSecret = getVar("auth0_client_secret", "AUTH0_CLIENT_SECRET", &diags)
	}
	return clientID, clientSecret, diags
}
