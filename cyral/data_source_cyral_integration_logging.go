package cyral

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ListIntegrationLogsResponse struct {
	Integrations []LoggingIntegration `json:"integrations"`
}

func (resp *ListIntegrationLogsResponse) WriteToSchema(d *schema.ResourceData) error {
	var integrationList []interface{}
	for _, integration := range resp.Integrations {
		// write in config scheme
		configScheme, err := getLoggingConfig(&integration)
		if err != nil {
			return err
		}

		integrationList = append(integrationList, map[string]interface{}{
			"id":                 integration.Id,
			"name":               integration.Name,
			"receive_audit_logs": integration.ReceiveAuditLogs,
			"config":             configScheme,
		})
	}
	if err := d.Set("integration_list", integrationList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	return nil
}

func dataSourceIntegrationLogsRead() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "IntegrationLogsDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			query := urlQuery(map[string]string{
				"type": d.Get("type").(string),
			})
			return fmt.Sprintf("https://%s/v1/integrations/logging%s", c.ControlPlane, query)
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &ListIntegrationLogsResponse{} },
	}
}

func dataSourceIntegrationLogging() *schema.Resource {
	rawSchema := getIntegrationLogsSchema()
	// all fields in integration_list are computed.
	// this function changes the schema to achieve this
	computedSchema := convertSchemaFieldsToComputed(rawSchema)
	log.Printf("[INFO] Computed schema: %v", computedSchema)
	return &schema.Resource{
		Description: "Retrieve and filter logging integrations.",
		ReadContext: ReadResource(dataSourceIntegrationLogsRead()),
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
			"integration_list": {
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
