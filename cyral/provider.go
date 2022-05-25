package cyral

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keycloak            = "keycloak"
	auth0               = "auth0"
	EnvVarClientID      = "CYRAL_TF_CLIENT_ID"
	EnvVarClientSecret  = "CYRAL_TF_CLIENT_SECRET"
	EnvVarCPURL         = "CYRAL_TF_CONTROL_PLANE"
	EnvVarTLSSkipVerify = "CYRAL_TF_TLS_SKIP_VERIFY"
)

func init() {
	schema.ResourceDescriptionBuilder = func(s *schema.Resource) string {
		desc := s.Description
		if s.DeprecationMessage != "" {
			desc = fmt.Sprintf("**Deprecated.** %s", s.DeprecationMessage)
		}
		return strings.TrimSpace(desc)
	}
}

// Provider defines and initializes the Cyral provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_provider": {
				Description: "Auth0-based control planes are no longer supported. Use `keycloak` " +
					"or remove the variable declaration",
				Type:       schema.TypeString,
				Optional:   true,
				Default:    keycloak,
				Deprecated: "Auth0-based control planes are no longer supported.",
			},
			"auth0_audience": {
				Description: "Auth0 audience.",
				Type:        schema.TypeString,
				Optional:    true,
				RequiredWith: []string{
					"auth0_domain",
				},
				Deprecated: "Auth0-based control planes are no longer supported.",
			},
			"auth0_domain": {
				Description: "Auth0 domain name.",
				Type:        schema.TypeString,
				Optional:    true,
				RequiredWith: []string{
					"auth0_audience",
				},
				Deprecated: "Auth0-based control planes are no longer supported.",
			},
			"auth0_client_id": {
				Description:   "Auth0 client id.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AUTH0_CLIENT_ID", nil),
				ConflictsWith: []string{"client_id"},
				Deprecated:    "Auth0-based control planes are no longer supported.",
			},
			"auth0_client_secret": {
				Description:   "Auth0 client secret.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("AUTH0_CLIENT_SECRET", nil),
				ConflictsWith: []string{"client_secret"},
				Deprecated:    "Auth0-based control planes are no longer supported.",
			},
			"client_id": {
				Description: "Client id used to authenticate against the control plane. Can be ommited and " +
					"declared using the environment variable `CYRAL_TF_CLIENT_ID`.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"auth0_client_id"},
				DefaultFunc:   schema.EnvDefaultFunc(EnvVarClientID, nil),
			},
			"client_secret": {
				Description: "Client secret used to authenticate against the control plane. Can be ommited and " +
					"declared using the environment variable `CYRAL_TF_CLIENT_SECRET`.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"auth0_client_secret"},
				DefaultFunc:   schema.EnvDefaultFunc(EnvVarClientSecret, nil),
			},
			"control_plane": {
				Description: "Control plane host and API port (ex: `some-cp.cyral.com:8000`)",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(EnvVarCPURL, nil),
			},
			"tls_skip_verify": {
				Type: schema.TypeBool,
				Description: "Specifies if the client will verify the TLS server certificate " +
					"used by the control plane. If set to `true`, the client will not verify " +
					"the server certificate, hence, it will allow insecure connections to be " +
					"established. This should be set only for testing and is not recommended " +
					"to be used in production environments. Can be set through the " +
					"`CYRAL_TF_TLS_SKIP_VERIFY` environment variable. Defaults to `false`.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(EnvVarTLSSkipVerify, nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cyral_saml_certificate":     dataSourceSAMLCertificate(),
			"cyral_saml_configuration":   dataSourceSAMLConfiguration(),
			"cyral_sidecar_bound_ports":  dataSourceSidecarBoundPorts(),
			"cyral_sidecar_cft_template": dataSourceSidecarCftTemplate(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cyral_datamap":                     resourceDatamap(),
			"cyral_identity_map":                resourceRepositoryIdentityMap("Use `cyral_repository_identity_map` instead."),
			"cyral_integration_datadog":         resourceIntegrationDatadog(),
			"cyral_integration_elk":             resourceIntegrationELK(),
			"cyral_integration_hc_vault":        resourceIntegrationHCVault(),
			"cyral_integration_logstash":        resourceIntegrationLogstash(),
			"cyral_integration_looker":          resourceIntegrationLooker(),
			"cyral_integration_microsoft_teams": resourceIntegrationMsTeams(),
			"cyral_integration_okta":            resourceIntegrationOkta(),
			"cyral_integration_pager_duty":      resourceIntegrationPagerDuty(),
			"cyral_integration_slack_alerts":    resourceIntegrationSlackAlerts(),
			"cyral_integration_splunk":          resourceIntegrationSplunk(),
			"cyral_integration_idp_aad":         resourceIntegrationIdP("aad", ""),
			"cyral_integration_idp_adfs":        resourceIntegrationIdP("adfs-2016", ""),
			"cyral_integration_idp_forgerock":   resourceIntegrationIdP("forgerock", ""),
			"cyral_integration_idp_gsuite":      resourceIntegrationIdP("gsuite", ""),
			"cyral_integration_idp_okta":        resourceIntegrationIdP("okta", ""),
			"cyral_integration_idp_ping_one":    resourceIntegrationIdP("pingone", ""),
			"cyral_integration_sso_aad":         resourceIntegrationIdP("aad", "Use 'cyral_integration_idp_aad' instead"),
			"cyral_integration_sso_adfs":        resourceIntegrationIdP("adfs-2016", "Use 'cyral_integration_idp_adfs' instead"),
			"cyral_integration_sso_forgerock":   resourceIntegrationIdP("forgerock", "Use 'cyral_integration_idp_forgerock' instead"),
			"cyral_integration_sso_gsuite":      resourceIntegrationIdP("gsuite", "Use 'cyral_integration_idp_gsuite' instead"),
			"cyral_integration_sso_okta":        resourceIntegrationIdP("okta", "Use 'cyral_integration_idp_okta' instead"),
			"cyral_integration_sso_ping_one":    resourceIntegrationIdP("pingone", "Use 'cyral_integration_idp_ping_one' instead"),
			"cyral_integration_sumo_logic":      resourceIntegrationSumoLogic(),
			"cyral_policy":                      resourcePolicy(),
			"cyral_policy_rule":                 resourcePolicyRule(),
			"cyral_repository":                  resourceRepository(),
			"cyral_repository_binding":          resourceRepositoryBinding(),
			"cyral_repository_conf_analysis":    resourceRepositoryConfAnalysis(),
			"cyral_repository_conf_auth":        resourceRepositoryConfAuth(),
			"cyral_repository_identity_map":     resourceRepositoryIdentityMap(""),
			"cyral_repository_local_account":    resourceRepositoryLocalAccount(),
			"cyral_role":                        resourceRole(),
			"cyral_role_sso_groups":             resourceRoleSSOGroups(),
			"cyral_sidecar":                     resourceSidecar(),
			"cyral_sidecar_credentials":         resourceSidecarCredentials(),
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
	tlsSkipVerify := d.Get("tls_skip_verify").(bool)

	log.Printf("[DEBUG] auth0Domain: %s ; auth0Audience: %s ; controlPlane: %s ; tlsSkipVerify: %t",
		auth0Domain, clientSecret, controlPlane, tlsSkipVerify)

	c, err := client.NewClient(clientID, clientSecret, auth0Domain, auth0Audience,
		controlPlane, keycloakProvider, tlsSkipVerify)
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
