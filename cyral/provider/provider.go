package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
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
	ctx := context.Background()
	tflog.Debug(ctx, "Init getDataSourceMap")
	schemaMap := map[string]*schema.Resource{}
	for _, p := range ps {
		tflog.Debug(ctx, fmt.Sprintf("Looking for datasources in package `%s`", p.Name()))
		for _, v := range p.Schemas() {
			if v.Type == core.DataSourceSchemaType {
				tflog.Debug(ctx, fmt.Sprintf("Registering datasources `%s`", v.Name))
				schemaMap[v.Name] = v.Schema()
			}
		}
	}

	schemaMap["cyral_integration_idp"] = deprecated.DataSourceIntegrationIdP()
	schemaMap["cyral_saml_configuration"] = deprecated.DataSourceSAMLConfiguration()
	schemaMap["cyral_sidecar_cft_template"] = deprecated.DataSourceSidecarCftTemplate()
	schemaMap["cyral_sidecar_instance_ids"] = deprecated.DataSourceSidecarInstanceIDs()

	schemaMap["cyral_sidecar_bound_ports"] = sidecar.DataSourceSidecarBoundPorts()
	schemaMap["cyral_sidecar_id"] = sidecar.DataSourceSidecarID()

	tflog.Debug(ctx, "End getDataSourceMap")

	return schemaMap
}

func getResourceMap(ps []core.PackageSchema) map[string]*schema.Resource {
	ctx := context.Background()
	tflog.Debug(ctx, "Init getResourceMap")
	var idpDeprecationMessage = "Use resource and data source `cyral_integration_idp_saml` instead."
	schemaMap := map[string]*schema.Resource{}
	for _, p := range ps {
		tflog.Debug(ctx, fmt.Sprintf("Looking for resources in package `%s`", p.Name()))
		for _, v := range p.Schemas() {
			if v.Type == core.ResourceSchemaType {
				tflog.Debug(ctx, fmt.Sprintf("Registering resources `%s`", v.Name))
				schemaMap[v.Name] = v.Schema()
			}
		}
	}

	// TODO Remove all the following resources in the next major version.
	schemaMap["cyral_integration_datadog"] = deprecated.ResourceIntegrationDatadog()
	schemaMap["cyral_integration_elk"] = deprecated.ResourceIntegrationELK()
	schemaMap["cyral_integration_logstash"] = deprecated.ResourceIntegrationLogstash()
	schemaMap["cyral_integration_looker"] = deprecated.ResourceIntegrationLooker()
	schemaMap["cyral_integration_splunk"] = deprecated.ResourceIntegrationSplunk()
	schemaMap["cyral_integration_idp_aad"] = deprecated.ResourceIntegrationIdP("aad", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_adfs"] = deprecated.ResourceIntegrationIdP("adfs-2016", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_forgerock"] = deprecated.ResourceIntegrationIdP("forgerock", "")
	schemaMap["cyral_integration_idp_gsuite"] = deprecated.ResourceIntegrationIdP("gsuite", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_okta"] = deprecated.ResourceIntegrationIdP("okta", idpDeprecationMessage)
	schemaMap["cyral_integration_idp_ping_one"] = deprecated.ResourceIntegrationIdP("pingone", idpDeprecationMessage)
	schemaMap["cyral_integration_sumo_logic"] = deprecated.ResourceIntegrationSumoLogic()

	tflog.Debug(ctx, "End getResourceMap")

	return schemaMap
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tflog.Debug(ctx, "Init providerConfigure")

	clientID, clientSecret, diags := getCredentials(d)
	if diags.HasError() {
		return nil, diags
	}

	controlPlane := d.Get("control_plane").(string)
	tlsSkipVerify := d.Get("tls_skip_verify").(bool)
	tflog.Debug(ctx, fmt.Sprintf("controlPlane: %s ; tlsSkipVerify: %t", controlPlane, tlsSkipVerify))

	c, err := client.New(clientID, clientSecret, controlPlane, tlsSkipVerify)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Cyral client",
			Detail:   err.Error(),
		})

		return nil, diags
	}
	tflog.Debug(ctx, "End providerConfigure")

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

var provider = Provider()

var ProviderFactories = map[string]func() (*schema.Provider, error){
	"cyral": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}
