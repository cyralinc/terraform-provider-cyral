package logging

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceIntegrationLogsRead() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "IntegrationLogsDataSourceRead",
		Type:         operationtype.Read,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			query := utils.UrlQuery(map[string]string{
				"type": d.Get("type").(string),
			})
			return fmt.Sprintf("https://%s/v1/integrations/logging%s", c.ControlPlane, query)
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &ListIntegrationLogsResponse{} },
	}
}

func dataSourceSchema() *schema.Resource {
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
