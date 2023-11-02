package pagerduty

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	ce "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIntegrationPagerDuty() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with PagerDuty](https://cyral.com/docs/integrations/incident-response/pagerduty/#in-cyral).",
		CreateContext: core.CreateResource(
			ce.ConfExtensionIntegrationCreate(ce.PagerDutyTemplateType),
			ce.ConfExtensionIntegrationRead(ce.PagerDutyTemplateType)),
		ReadContext: core.ReadResource(ce.ConfExtensionIntegrationRead(ce.PagerDutyTemplateType)),
		UpdateContext: core.UpdateResource(
			ce.ConfExtensionIntegrationUpdate(ce.PagerDutyTemplateType),
			ce.ConfExtensionIntegrationRead(ce.PagerDutyTemplateType)),
		DeleteContext: core.DeleteResource(ce.ConfExtensionIntegrationDelete(ce.PagerDutyTemplateType)),

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
