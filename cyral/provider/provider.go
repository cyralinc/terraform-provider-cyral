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
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/awsiam"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension/mfaduo"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension/pagerduty"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/hcvault"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/logging"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/slack"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/teams"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/permission"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/policy"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/policy/rule"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/regopolicy"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessgateway"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessrules"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/binding"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confanalysis"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confauth"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/network"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/useraccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/role"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/samlcertificate"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/samlconfiguration"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/serviceaccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/credentials"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/health"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/instance"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/listener"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/systeminfo"
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

	schemaMap["cyral_integration_idp"] = deprecated.DataSourceIntegrationIdP()
	schemaMap["cyral_integration_idp_saml"] = idpsaml.DataSourceIntegrationIdPSAML()
	schemaMap["cyral_integration_logging"] = logging.DataSourceIntegrationLogging()
	schemaMap["cyral_permission"] = permission.DataSourcePermission()
	schemaMap["cyral_repository"] = repository.DataSourceRepository()
	schemaMap["cyral_role"] = role.DataSourceRole()
	schemaMap["cyral_saml_certificate"] = samlcertificate.DataSourceSAMLCertificate()
	schemaMap["cyral_saml_configuration"] = samlconfiguration.DataSourceSAMLConfiguration()
	schemaMap["cyral_sidecar_bound_ports"] = sidecar.DataSourceSidecarBoundPorts()
	schemaMap["cyral_sidecar_cft_template"] = deprecated.DataSourceSidecarCftTemplate()
	schemaMap["cyral_sidecar_health"] = health.DataSourceSidecarHealth()
	schemaMap["cyral_sidecar_id"] = sidecar.DataSourceSidecarID()
	schemaMap["cyral_sidecar_instance_ids"] = deprecated.DataSourceSidecarInstanceIDs()
	schemaMap["cyral_sidecar_instance_stats"] = instance.DataSourceSidecarInstanceStats()
	schemaMap["cyral_sidecar_instance"] = instance.DataSourceSidecarInstance()
	schemaMap["cyral_sidecar_listener"] = listener.DataSourceSidecarListener()
	schemaMap["cyral_system_info"] = systeminfo.DataSourceSystemInfo()

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
	schemaMap["cyral_integration_aws_iam"] = awsiam.ResourceIntegrationAWSIAM()
	schemaMap["cyral_integration_datadog"] = deprecated.ResourceIntegrationDatadog()
	schemaMap["cyral_integration_mfa_duo"] = mfaduo.ResourceIntegrationMFADuo()
	schemaMap["cyral_integration_elk"] = deprecated.ResourceIntegrationELK()
	schemaMap["cyral_integration_hc_vault"] = hcvault.ResourceIntegrationHCVault()
	schemaMap["cyral_integration_logstash"] = deprecated.ResourceIntegrationLogstash()
	schemaMap["cyral_integration_looker"] = deprecated.ResourceIntegrationLooker()
	schemaMap["cyral_integration_microsoft_teams"] = teams.ResourceIntegrationMsTeams()
	schemaMap["cyral_integration_pager_duty"] = pagerduty.ResourceIntegrationPagerDuty()
	schemaMap["cyral_integration_slack_alerts"] = slack.ResourceIntegrationSlackAlerts()
	schemaMap["cyral_integration_splunk"] = deprecated.ResourceIntegrationSplunk()
	schemaMap["cyral_integration_idp_aad"] = deprecated.ResourceIntegrationIdP("aad", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_adfs"] = deprecated.ResourceIntegrationIdP("adfs-2016", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_forgerock"] = deprecated.ResourceIntegrationIdP("forgerock", "")
	schemaMap["cyral_integration_idp_gsuite"] = deprecated.ResourceIntegrationIdP("gsuite", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_okta"] = deprecated.ResourceIntegrationIdP("okta", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_ping_one"] = deprecated.ResourceIntegrationIdP("pingone", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_saml"] = idpsaml.ResourceIntegrationIdPSAML()
	schemaMap["cyral_integration_idp_saml_draft"] = idpsaml.ResourceIntegrationIdPSAMLDraft()
	schemaMap["cyral_integration_sumo_logic"] = deprecated.ResourceIntegrationSumoLogic()
	schemaMap["cyral_integration_logging"] = logging.ResourceIntegrationLogging()
	schemaMap["cyral_policy"] = policy.ResourcePolicy()
	schemaMap["cyral_policy_rule"] = rule.ResourcePolicyRule()
	schemaMap["cyral_rego_policy_instance"] = regopolicy.ResourceRegoPolicyInstance()
	schemaMap["cyral_repository"] = repository.ResourceRepository()
	schemaMap["cyral_repository_access_rules"] = accessrules.ResourceRepositoryAccessRules()
	schemaMap["cyral_repository_access_gateway"] = accessgateway.ResourceRepositoryAccessGateway()
	schemaMap["cyral_repository_binding"] = binding.ResourceRepositoryBinding()
	schemaMap["cyral_repository_conf_auth"] = confauth.ResourceRepositoryConfAuth()
	schemaMap["cyral_repository_conf_analysis"] = confanalysis.ResourceRepositoryConfAnalysis()
	schemaMap["cyral_repository_network_access_policy"] = network.ResourceRepositoryNetworkAccessPolicy()
	schemaMap["cyral_repository_user_account"] = useraccount.ResourceRepositoryUserAccount()
	schemaMap["cyral_role"] = role.ResourceRole()
	schemaMap["cyral_role_sso_groups"] = role.ResourceRoleSSOGroups()
	schemaMap["cyral_service_account"] = serviceaccount.ResourceServiceAccount()
	schemaMap["cyral_sidecar"] = sidecar.ResourceSidecar()
	schemaMap["cyral_sidecar_credentials"] = credentials.ResourceSidecarCredentials()
	schemaMap["cyral_sidecar_listener"] = listener.ResourceSidecarListener()

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

var ProviderFactories = map[string]func() (*schema.Provider, error){
	"cyral": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}
