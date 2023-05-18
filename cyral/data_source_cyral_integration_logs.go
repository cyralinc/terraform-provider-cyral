package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ListIntegrationLogsRequest struct {
	IntegrationType string `json:"type"`
}

type ListIntegrationLogsResponse struct {
	Integrations []IntegrationLogConfig `json:"integrations"`
}

func (resp *ListIntegrationLogsResponse) WriteToSchema(d *schema.ResourceData) error {
	var integrationList []interface{}
	for _, integration := range resp.Integrations {
		// if err := d.Set("integration_id", integration.Id); err != nil {
		// 	return fmt.Errorf("error setting 'integration_id': %w", err)
		// }

		if err := d.Set("name", integration.Name); err != nil {
			return fmt.Errorf("error setting 'name': %w", err)
		}
		if err := d.Set("enable_audit_logs", integration.EnableAuditLogs); err != nil {
			return fmt.Errorf("error setting 'enable_audit_logs': %w", err)
		}

		// write in config scheme
		configScheme, err := writeConfigScheme(&integration)
		if err != nil {
			return err
		}

		integrationList = append(integrationList, map[string]interface{}{
			"integration_id":    integration.Id,
			"name":              integration.Name,
			"enable_audit_logs": integration.EnableAuditLogs,
			"config_scheme":     configScheme,
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

func dataSourceIntegrationLogs() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter integrations.",
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
				}, false),
			},
			"integration_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: getIntegrationLogsSchema(),
				},
				Computed:    true,
				Description: "List of existing integration configs for the given filter criteria.",
			},
		},
	}
}
