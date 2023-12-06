package awsiam

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	AWSIAMIntegrationNameKey        = "name"
	AWSIAMIntegratioNDescriptionKey = "description"
	AWSIAMIntegrationARNsKey        = "role_arns"
)

type AWSIAMIntegrationWrapper struct {
	Integration *AWSIAMIntegration `json:"iamIntegration"`
}

type AWSIAMIntegration struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IAMRoleARNs []string `json:"iamRoleARNs"`
}

func (wrapper *AWSIAMIntegrationWrapper) WriteToSchema(d *schema.ResourceData) error {
	integration := wrapper.Integration

	d.SetId(integration.ID)

	if err := d.Set(AWSIAMIntegrationNameKey, integration.Name); err != nil {
		return fmt.Errorf("error setting '%s': %w", AWSIAMIntegrationNameKey, err)
	}

	if err := d.Set(AWSIAMIntegratioNDescriptionKey, integration.Description); err != nil {
		return fmt.Errorf("error setting '%s': %w", AWSIAMIntegratioNDescriptionKey, err)
	}

	if err := d.Set(AWSIAMIntegrationARNsKey, integration.IAMRoleARNs); err != nil {
		return fmt.Errorf("error setting '%s': %w", AWSIAMIntegrationARNsKey, err)
	}
	return nil
}

func (wrapper *AWSIAMIntegrationWrapper) ReadFromSchema(d *schema.ResourceData) error {
	wrapper.Integration = &AWSIAMIntegration{}

	wrapper.Integration.Name = d.Get(AWSIAMIntegrationNameKey).(string)
	wrapper.Integration.Description = d.Get(AWSIAMIntegratioNDescriptionKey).(string)

	arns := d.Get(AWSIAMIntegrationARNsKey).([]interface{})
	stringARNs := make([]string, 0, len(arns))
	for _, arn := range arns {
		stringARNs = append(stringARNs, arn.(string))
	}

	wrapper.Integration.IAMRoleARNs = stringARNs
	return nil
}

func ResourceIntegrationAWSIAM() *schema.Resource {
	contextHandler := core.DefaultContextHandler{
		ResourceName:        "AWS IAM Integration",
		ResourceType:        resourcetype.Resource,
		SchemaReaderFactory: func() core.SchemaReader { return &AWSIAMIntegrationWrapper{} },
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &AWSIAMIntegrationWrapper{} },
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
