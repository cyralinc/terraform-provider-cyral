package awsiam

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSchema() *schema.Resource {
	contextHandler := core.DefaultContextHandler{
		ResourceName:                 "AWS IAM Integration",
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &AWSIAMIntegrationWrapper{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &AWSIAMIntegrationWrapper{} },
		BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/aws/iam", c.ControlPlane)
		},
	}
	return &schema.Resource{
		Description:   "Authenticate users based on AWS IAM credentials.",
		CreateContext: contextHandler.CreateContext(),
		ReadContext:   contextHandler.ReadContext(),
		UpdateContext: contextHandler.UpdateContext(),
		DeleteContext: contextHandler.DeleteContext(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			utils.IDKey: {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			AWSIAMIntegrationNameKey: {
				Description: "The name of this AWS IAM Authentication integration.",
				Required:    true,
				Type:        schema.TypeString,
			},

			AWSIAMIntegratioNDescriptionKey: {
				Description: "Optional description of this integration.",
				Optional:    true,
				Type:        schema.TypeString,
			},

			AWSIAMIntegrationARNsKey: {
				Description: "List of role ARNs which will be used for authentication.",
				Required:    true,
				MinItems:    1,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
