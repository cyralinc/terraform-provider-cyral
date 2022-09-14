package cyral

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationPagerDuty() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with PagerDuty](https://cyral.com/docs/integrations/incident-response/pagerduty/#in-cyral).",
		CreateContext: CreateResource(
			ConfExtensionIntegrationCreate(pagerDutyTemplateType),
			ConfExtensionIntegrationRead(pagerDutyTemplateType)),
		ReadContext: ReadResource(ConfExtensionIntegrationRead(pagerDutyTemplateType)),
		UpdateContext: UpdateResource(
			ConfExtensionIntegrationUpdate(pagerDutyTemplateType),
			ConfExtensionIntegrationRead(pagerDutyTemplateType)),
		DeleteContext: DeleteResource(ConfExtensionIntegrationDelete(pagerDutyTemplateType)),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"api_token": {
				Description: "API token for the PagerDuty integration.",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
