package cyral

import (
	"context"
	"fmt"
	"log"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keycloak           = "keycloak"
	auth0              = "auth0"
	EnvVarClientID     = "CYRAL_TF_CLIENT_ID"
	EnvVarClientSecret = "CYRAL_TF_CLIENT_SECRET"
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
				Optional: true,
				RequiredWith: []string{
					"auth0_domain",
				},
			},
			"auth0_domain": {
				Type:     schema.TypeString,
				Optional: true,
				RequiredWith: []string{
					"auth0_audience",
				},
			},
			"auth0_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AUTH0_CLIENT_ID", nil),
				ConflictsWith: []string{"client_id"},
				Deprecated: fmt.Sprintf("use provider variable 'client_id' or environment variable "+
					"'%s' instead of 'auth0_client_id'", EnvVarClientID),
			},
			"auth0_client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AUTH0_CLIENT_SECRET", nil),
				ConflictsWith: []string{"client_secret"},
				Deprecated: fmt.Sprintf("use provider variable 'client_secret' or environment variable "+
					"'%s' instead of 'auth0_client_secret'", EnvVarClientSecret),
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"auth0_client_id"},
				DefaultFunc:   schema.EnvDefaultFunc(EnvVarClientID, nil),
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"auth0_client_secret"},
				DefaultFunc:   schema.EnvDefaultFunc(EnvVarClientSecret, nil),
			},
			"control_plane": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cyral_datamap":                     resourceDatamap(),
			"cyral_integration_datadog":         resourceIntegrationDatadog(),
			"cyral_integration_elk":             resourceIntegrationELK(),
			"cyral_integration_logstash":        resourceIntegrationLogstash(),
			"cyral_integration_looker":          resourceIntegrationLooker(),
			"cyral_integration_microsoft_teams": resourceIntegrationTeams(),
			"cyral_integration_splunk":          resourceIntegrationSplunk(),
			"cyral_integration_sumo_logic":      resourceIntegrationSumoLogic(),
			"cyral_policy":                      resourcePolicy(),
			"cyral_policy_rule":                 resourcePolicyRule(),
			"cyral_repository":                  resourceRepository(),
			"cyral_repository_binding":          resourceRepositoryBinding(),
			"cyral_sidecar":                     resourceSidecar(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[DEBUG] Init providerConfigure")
	keycloakProvider := d.Get("auth_provider").(string) == keycloak

	log.Printf("[DEBUG] keycloakProvider: %v", keycloakProvider)
	clientID, clientSecret, diags := getCredentials(d, keycloakProvider)

	if clientID == "" || clientSecret == "" {
		return nil, diags
	}
	log.Printf("[DEBUG] clientID: %s ; clientSecret: %s", clientID, clientSecret)

	auth0Domain := d.Get("auth0_domain").(string)
	auth0Audience := d.Get("auth0_audience").(string)
	controlPlane := d.Get("control_plane").(string)

	log.Printf("[DEBUG] auth0Domain: %s ; auth0Audience: %s ; controlPlane: %s",
		auth0Domain, clientSecret, controlPlane)

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
	log.Printf("[DEBUG] End providerConfigure")

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

	clientID = getVar("client_id", EnvVarClientID, &diags)
	clientSecret = getVar("client_secret", EnvVarClientSecret, &diags)

	// Backwards compatibility code to allow users to migrate to new variables and see
	// a deprecation warning. The code below must be removed in next versions.
	if !keycloakProvider && clientID == "" && clientSecret == "" {
		diags = nil
		clientID = getVar("auth0_client_id", "AUTH0_CLIENT_ID", &diags)
		clientSecret = getVar("auth0_client_secret", "AUTH0_CLIENT_SECRET", &diags)
	}
	return clientID, clientSecret, diags
}

func createError(summary, detail string) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   detail,
	})

	return diags
}
