package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
)

func init() {
	schema.ResourceDescriptionBuilder = func(s *schema.Resource) string {
		desc := s.Description
		if s.DeprecationMessage != "" {
			desc = fmt.Sprintf("\n~> **DEPRECATED** %s", s.DeprecationMessage)
		}
		return strings.TrimSpace(desc)
	}
}

// Provider defines and initializes the Cyral provider
func Provider() *schema.Provider {
	ps := packagesSchemas()
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Description: "Client id used to authenticate against the control plane. Can be ommited and " +
					"declared using the environment variable `CYRAL_TF_CLIENT_ID`.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvVarClientID, nil),
			},
			"client_secret": {
				Description: "Client secret used to authenticate against the control plane. Can be ommited and " +
					"declared using the environment variable `CYRAL_TF_CLIENT_SECRET`.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvVarClientSecret, nil),
			},
			"control_plane": {
				Description: "Control plane host and API port (ex: `tenant.app.cyral.com`)",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvVarCPURL, nil),
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
				DefaultFunc: schema.EnvDefaultFunc(client.EnvVarTLSSkipVerify, nil),
			},
		},
		DataSourcesMap:       getDataSourceMap(ps),
		ResourcesMap:         getResourceMap(ps),
		ConfigureContextFunc: providerConfigure,
	}
}

func getDataSourceMap(ps []core.PackageSchema) map[string]*schema.Resource {
	log.Printf("[DEBUG] Init getDataSourceMap")
	schemaMap := map[string]*schema.Resource{}
	for _, p := range ps {
		log.Printf("[DEBUG] Looking for datasources in package `%s`", p.Name())
		for _, v := range p.Schemas() {
			if v.Type == core.DataSourceSchemaType {
				log.Printf("[DEBUG] Registering datasources `%s`", v.Name)
				schemaMap[v.Name] = v.Schema()
			}
		}
	}
	schemaMap["cyral_integration_idp"] = internal.DataSourceIntegrationIdP()
	schemaMap["cyral_integration_idp_saml"] = internal.DataSourceIntegrationIdPSAML()
	schemaMap["cyral_integration_logging"] = internal.DataSourceIntegrationLogging()
	schemaMap["cyral_repository"] = internal.DataSourceRepository()
	schemaMap["cyral_role"] = internal.DataSourceRole()
	schemaMap["cyral_saml_certificate"] = internal.DataSourceSAMLCertificate()
	schemaMap["cyral_saml_configuration"] = internal.DataSourceSAMLConfiguration()
	schemaMap["cyral_sidecar_bound_ports"] = internal.DataSourceSidecarBoundPorts()
	schemaMap["cyral_sidecar_cft_template"] = internal.DataSourceSidecarCftTemplate()
	schemaMap["cyral_sidecar_id"] = internal.DataSourceSidecarID()
	schemaMap["cyral_sidecar_instance_ids"] = internal.DataSourceSidecarInstanceIDs()
	schemaMap["cyral_sidecar_listener"] = internal.DataSourceSidecarListener()

	log.Printf("[DEBUG] end getDataSourceMap")

	return schemaMap
}

func getResourceMap(ps []core.PackageSchema) map[string]*schema.Resource {
	log.Printf("[DEBUG] Init getResourceMap")
	var idpDeprecationMessage = "Use resource and data source `cyral_integration_idp_saml` instead."
	schemaMap := map[string]*schema.Resource{}
	for _, p := range ps {
		log.Printf("[DEBUG] Looking for resources in package `%s`", p.Name())
		for _, v := range p.Schemas() {
			if v.Type == core.ResourceSchemaType {
				log.Printf("[DEBUG] Registering resources `%s`", v.Name)
				schemaMap[v.Name] = v.Schema()
			}
		}
	}

	// // TODO Once the resources are migrated to the new SchemaRegister
	// // abstraction, these calls from provider to resource will be removed.
	schemaMap["cyral_integration_datadog"] = internal.ResourceIntegrationDatadog()
	schemaMap["cyral_integration_mfa_duo"] = internal.ResourceIntegrationMFADuo()
	schemaMap["cyral_integration_elk"] = internal.ResourceIntegrationELK()
	schemaMap["cyral_integration_hc_vault"] = internal.ResourceIntegrationHCVault()
	schemaMap["cyral_integration_logstash"] = internal.ResourceIntegrationLogstash()
	schemaMap["cyral_integration_looker"] = internal.ResourceIntegrationLooker()
	schemaMap["cyral_integration_microsoft_teams"] = internal.ResourceIntegrationMsTeams()
	schemaMap["cyral_integration_pager_duty"] = internal.ResourceIntegrationPagerDuty()
	schemaMap["cyral_integration_slack_alerts"] = internal.ResourceIntegrationSlackAlerts()
	schemaMap["cyral_integration_splunk"] = internal.ResourceIntegrationSplunk()
	schemaMap["cyral_integration_idp_aad"] = internal.ResourceIntegrationIdP("aad", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_adfs"] = internal.ResourceIntegrationIdP("adfs-2016", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_forgerock"] = internal.ResourceIntegrationIdP("forgerock", "")
	schemaMap["cyral_integration_idp_gsuite"] = internal.ResourceIntegrationIdP("gsuite", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_okta"] = internal.ResourceIntegrationIdP("okta", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_ping_one"] = internal.ResourceIntegrationIdP("pingone", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_saml"] = internal.ResourceIntegrationIdPSAML()
	schemaMap["cyral_integration_idp_saml_draft"] = internal.ResourceIntegrationIdPSAMLDraft()
	schemaMap["cyral_integration_sumo_logic"] = internal.ResourceIntegrationSumoLogic()
	schemaMap["cyral_integration_logging"] = internal.ResourceIntegrationLogging()
	schemaMap["cyral_policy"] = internal.ResourcePolicy()
	schemaMap["cyral_policy_rule"] = internal.ResourcePolicyRule()
	schemaMap["cyral_rego_policy_instance"] = internal.ResourceRegoPolicyInstance()
	schemaMap["cyral_repository"] = internal.ResourceRepository()
	schemaMap["cyral_repository_access_rules"] = internal.ResourceRepositoryAccessRules()
	schemaMap["cyral_repository_access_gateway"] = internal.ResourceRepositoryAccessGateway()
	schemaMap["cyral_repository_binding"] = internal.ResourceRepositoryBinding()
	schemaMap["cyral_repository_conf_auth"] = internal.ResourceRepositoryConfAuth()
	schemaMap["cyral_repository_conf_analysis"] = internal.ResourceRepositoryConfAnalysis()
	schemaMap["cyral_repository_network_access_policy"] = internal.ResourceRepositoryNetworkAccessPolicy()
	schemaMap["cyral_repository_user_account"] = internal.ResourceRepositoryUserAccount()
	schemaMap["cyral_role"] = internal.ResourceRole()
	schemaMap["cyral_role_sso_groups"] = internal.ResourceRoleSSOGroups()
	schemaMap["cyral_sidecar"] = internal.ResourceSidecar()
	schemaMap["cyral_sidecar_credentials"] = internal.ResourceSidecarCredentials()
	schemaMap["cyral_sidecar_listener"] = internal.ResourceSidecarListener()

	log.Printf("[DEBUG] End getResourceMap")

	return schemaMap
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

	c, err := client.New(clientID, clientSecret, controlPlane, tlsSkipVerify)
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

	clientID = getVar("client_id", client.EnvVarClientID, &diags)
	clientSecret = getVar("client_secret", client.EnvVarClientSecret, &diags)

	return clientID, clientSecret, diags
}
