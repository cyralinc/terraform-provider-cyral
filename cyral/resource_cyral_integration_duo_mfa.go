package cyral

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationDuoMFA() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Duo MFA](https://cyral.com/docs/mfa/duo).",
		CreateContext: CreateResource(
			ConfExtensionIntegrationCreate(duoMFATemplateType),
			ConfExtensionIntegrationRead(duoMFATemplateType)),
		ReadContext: ReadResource(ConfExtensionIntegrationRead(duoMFATemplateType)),
		UpdateContext: UpdateResource(
			ConfExtensionIntegrationUpdate(duoMFATemplateType),
			ConfExtensionIntegrationRead(duoMFATemplateType)),
		DeleteContext: DeleteResource(ConfExtensionIntegrationDelete(duoMFATemplateType)),

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
			"integration_key": {
				Description: "API token for the duo MFA integration.",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
			"secret_key": {
				Description: "API token for the duo MFA integration.",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
			"api_hostname": {
				Description: "API token for the duo MFA integration.",
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
