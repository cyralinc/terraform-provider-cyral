package mfaduo

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	ce "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIntegrationMFADuo() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [integration with Duo MFA](https://cyral.com/docs/mfa/duo).",
		CreateContext: core.CreateResource(
			ce.ConfExtensionIntegrationCreate(ce.DuoMFATemplateType),
			ce.ConfExtensionIntegrationRead(ce.DuoMFATemplateType)),
		ReadContext: core.ReadResource(ce.ConfExtensionIntegrationRead(ce.DuoMFATemplateType)),
		UpdateContext: core.UpdateResource(
			ce.ConfExtensionIntegrationUpdate(ce.DuoMFATemplateType),
			ce.ConfExtensionIntegrationRead(ce.DuoMFATemplateType)),
		DeleteContext: core.DeleteResource(ce.ConfExtensionIntegrationDelete(ce.DuoMFATemplateType)),

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
