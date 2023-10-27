package internal

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIntegrationPagerDuty() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with PagerDuty](https://cyral.com/docs/integrations/incident-response/pagerduty/#in-cyral).",
		CreateContext: core.CreateResource(
			ConfExtensionIntegrationCreate(PagerDutyTemplateType),
			ConfExtensionIntegrationRead(PagerDutyTemplateType)),
		ReadContext: core.ReadResource(ConfExtensionIntegrationRead(PagerDutyTemplateType)),
		UpdateContext: core.UpdateResource(
			ConfExtensionIntegrationUpdate(PagerDutyTemplateType),
			ConfExtensionIntegrationRead(PagerDutyTemplateType)),
		DeleteContext: core.DeleteResource(ConfExtensionIntegrationDelete(PagerDutyTemplateType)),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Integration display name that will be used in the control plane.",
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
