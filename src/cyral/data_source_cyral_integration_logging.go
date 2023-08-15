package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/src/client"
	"github.com/cyralinc/terraform-provider-cyral/src/core"
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ListIntegrationLogsResponse struct {
	Integrations []LoggingIntegration `json:"integrations"`
}

func (resp *ListIntegrationLogsResponse) WriteToSchema(d *schema.ResourceData) error {
	integrations := make([]interface{}, len(resp.Integrations))
	for i, integration := range resp.Integrations {
		// write in config scheme
		configType, config, err := getLoggingConfig(&integration)
		if err != nil {
			return err
		}
		integrations[i] = map[string]interface{}{
			"id":                 integration.Id,
			"name":               integration.Name,
			"receive_audit_logs": integration.ReceiveAuditLogs,
			configType:           config,
		}
	}
	if err := d.Set("integrations", integrations); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	return nil
}

func dataSourceIntegrationLogsRead() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		Name:       "IntegrationLogsDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			query := utils.UrlQuery(map[string]string{
				"type": d.Get("type").(string),
			})
			return fmt.Sprintf("https://%s/v1/integrations/logging%s", c.ControlPlane, query)
		},
		NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &ListIntegrationLogsResponse{} },
	}
}

func dataSourceIntegrationLogging() *schema.Resource {
	rawSchema := getIntegrationLogsSchema()
	// all fields in integrations are computed.
	// this function changes the schema to achieve this
	computedSchema := utils.ConvertSchemaFieldsToComputed(rawSchema)
	return &schema.Resource{
		Description: "Retrieve and filter logging integrations.",
		ReadContext: core.ReadResource(dataSourceIntegrationLogsRead()),
		Schema: map[string]*schema.Schema{
			"type": {
				Description: "The type of logging integration config to filter by.",
				Optional:    true,
				Type:        schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"CLOUDWATCH",
					"DATADOG",
					"ELK",
					"SPLUNK",
					"SUMOLOGIC",
					"FLUENTBIT",
				}, false),
				Default: "ANY",
			},
			"integrations": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: computedSchema,
				},
				Computed:    true,
				Description: "List of existing integration configs for the given filter criteria.",
			},
		},
	}
}
