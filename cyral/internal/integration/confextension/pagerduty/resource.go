package pagerduty

import (
	ce "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [integration with PagerDuty](https://cyral.com/docs/integrations/incident-response/pagerduty/#in-cyral).",
		CreateContext: ce.CreateResource(resourceName, ce.PagerDutyTemplateType),
		ReadContext:   ce.ReadResource(resourceName, ce.PagerDutyTemplateType),
		UpdateContext: ce.UpdateResource(resourceName, ce.PagerDutyTemplateType),
		DeleteContext: ce.DeleteResource(resourceName, ce.PagerDutyTemplateType),

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
