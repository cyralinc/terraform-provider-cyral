package cyral

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

const (
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
			"client_id": {
				Description: "Client id used to authenticate against the control plane. Can be ommited and " +
					"declared using the environment variable `CYRAL_TF_CLIENT_ID`.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(EnvVarClientID, nil),
			},
			"client_secret": {
				Description: "Client secret used to authenticate against the control plane. Can be ommited and " +
					"declared using the environment variable `CYRAL_TF_CLIENT_SECRET`.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(EnvVarClientSecret, nil),
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
			"cyral_datalabel":            dataSourceDatalabel(),
			"cyral_integration_idp":      dataSourceIntegrationIdP(),
			"cyral_integration_idp_saml": dataSourceIntegrationIdPSAML(),
			"cyral_repository":           dataSourceRepository(),
			"cyral_role":                 dataSourceRole(),
			"cyral_saml_certificate":     dataSourceSAMLCertificate(),
			"cyral_saml_configuration":   dataSourceSAMLConfiguration(),
			"cyral_sidecar_bound_ports":  dataSourceSidecarBoundPorts(),
			"cyral_sidecar_cft_template": dataSourceSidecarCftTemplate(),
			"cyral_sidecar_id":           dataSourceSidecarID(),
			"cyral_sidecar_instance_ids": dataSourceSidecarInstanceIDs(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cyral_datalabel":                        resourceDatalabel(),
			"cyral_integration_datadog":              resourceIntegrationDatadog(),
			"cyral_integration_mfa_duo":              resourceIntegrationMFADuo(),
			"cyral_integration_elk":                  resourceIntegrationELK(),
			"cyral_integration_hc_vault":             resourceIntegrationHCVault(),
			"cyral_integration_logstash":             resourceIntegrationLogstash(),
			"cyral_integration_looker":               resourceIntegrationLooker(),
			"cyral_integration_microsoft_teams":      resourceIntegrationMsTeams(),
			"cyral_integration_pager_duty":           resourceIntegrationPagerDuty(),
			"cyral_integration_slack_alerts":         resourceIntegrationSlackAlerts(),
			"cyral_integration_splunk":               resourceIntegrationSplunk(),
			"cyral_integration_idp_aad":              resourceIntegrationIdP("aad"),
			"cyral_integration_idp_adfs":             resourceIntegrationIdP("adfs-2016"),
			"cyral_integration_idp_forgerock":        resourceIntegrationIdP("forgerock"),
			"cyral_integration_idp_gsuite":           resourceIntegrationIdP("gsuite"),
			"cyral_integration_idp_okta":             resourceIntegrationIdP("okta"),
			"cyral_integration_idp_ping_one":         resourceIntegrationIdP("pingone"),
			"cyral_integration_idp_saml":             resourceIntegrationIdPSAML(),
			"cyral_integration_idp_saml_draft":       resourceIntegrationIdPSAMLDraft(),
			"cyral_integration_sumo_logic":           resourceIntegrationSumoLogic(),
			"cyral_policy":                           resourcePolicy(),
			"cyral_policy_rule":                      resourcePolicyRule(),
			"cyral_repository":                       resourceRepository(),
			"cyral_repository_binding":               resourceRepositoryBinding(),
			"cyral_repository_conf_analysis":         resourceRepositoryConfAnalysis(),
			"cyral_repository_conf_auth":             resourceRepositoryConfAuth(),
			"cyral_repository_datamap":               resourceRepositoryDatamap(),
			"cyral_repository_user_account":          resourceRepositoryUserAccount(),
			"cyral_repository_network_access_policy": resourceRepositoryNetworkAccessPolicy(),
			"cyral_repository_access_rules":          resourceRepositoryAccessRules(),
			"cyral_role":                             resourceRole(),
			"cyral_role_sso_groups":                  resourceRoleSSOGroups(),
			"cyral_sidecar":                          resourceSidecar(),
			"cyral_sidecar_credentials":              resourceSidecarCredentials(),
			// The Sidecar Listener resource will be reenabled when the port-multiplexing
			// feature is completed. Jira: https://cyralinc.atlassian.net/browse/ENG-9398
			//"cyral_sidecar_listener":                 resourceSidecarListener(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[DEBUG] Init providerConfigure")

	clientID, clientSecret, diags := getCredentials(d)
	if diags.HasError() {
		return nil, diags
	}

	controlPlane := d.Get("control_plane").(string)
	tlsSkipVerify := d.Get("tls_skip_verify").(bool)
	log.Printf("[DEBUG] controlPlane: %s ; tlsSkipVerify: %t", controlPlane, tlsSkipVerify)

	c, err := client.NewClient(clientID, clientSecret, controlPlane, tlsSkipVerify)
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

func getCredentials(d *schema.ResourceData) (string, string, diag.Diagnostics) {
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

	return clientID, clientSecret, diags
}

func getProviderConfigFromEnv() (
	clientID string,
	clientSecret string,
	controlPlane string,
	tlsSkipVerify bool,
	err error,
) {
	clientID = os.Getenv(EnvVarClientID)
	clientSecret = os.Getenv(EnvVarClientSecret)
	controlPlane = os.Getenv(EnvVarCPURL)
	tlsSkipVerifyStr := os.Getenv(EnvVarTLSSkipVerify)
	if tlsSkipVerifyStr != "" {
		tlsSkipVerify, err = strconv.ParseBool(tlsSkipVerifyStr)
		if err != nil {
			return "", "", "", false, fmt.Errorf("invalid value for "+
				"env var %q: %w", EnvVarTLSSkipVerify, err)
		}
	}
	return
}

func newClientFromEnv() (*client.Client, error) {
	clientID, clientSecret, controlPlane, tlsSkipVerify, err :=
		getProviderConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("unable to create Cyral client: %w", err)
	}
	c, err := client.NewClient(clientID, clientSecret, controlPlane,
		tlsSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("unable to create Cyral client: %w", err)
	}
	return c, nil
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
