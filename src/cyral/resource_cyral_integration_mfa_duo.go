package cyral

import (
	"github.com/cyralinc/terraform-provider-cyral/src/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIntegrationMFADuo() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Duo MFA](https://cyral.com/docs/mfa/duo).",
		CreateContext: core.CreateResource(
			ConfExtensionIntegrationCreate(duoMFATemplateType),
			ConfExtensionIntegrationRead(duoMFATemplateType)),
		ReadContext: core.ReadResource(ConfExtensionIntegrationRead(duoMFATemplateType)),
		UpdateContext: core.UpdateResource(
			ConfExtensionIntegrationUpdate(duoMFATemplateType),
			ConfExtensionIntegrationRead(duoMFATemplateType)),
		DeleteContext: core.DeleteResource(ConfExtensionIntegrationDelete(duoMFATemplateType)),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"name": {
				Description:  "Integration display name that will be used in the control plane.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"integration_key": {
				Description:  "Integration key name obtained from Duo management console.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"secret_key": {
				Description:  "Secret key obtained from Duo management console.",
				Required:     true,
				Sensitive:    true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"api_hostname": {
				Description:  "API hostname obtained from Duo management console.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
