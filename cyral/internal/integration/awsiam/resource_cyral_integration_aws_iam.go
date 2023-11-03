package awsiam

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
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

var ReadAWSIAMIntegration = core.ResourceOperationConfig{
	Name:       "AWSIAMIntegrationRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/integrations/aws/iam/%s",
			c.ControlPlane,
			d.Id(),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
		return &AWSIAMIntegrationWrapper{}
	},
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "AWS IAM Integration"},
}

func ResourceIntegrationAWSIAM() *schema.Resource {
	return &schema.Resource{
		Description: "Authenticate users based on AWS IAM credentials.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "AWSIAMIntegrationCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/aws/iam", c.ControlPlane)
				},
				NewResourceData: func() core.ResourceData {
					return &AWSIAMIntegrationWrapper{}
				},
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
					return &core.IDBasedResponse{}
				},
			},
			ReadAWSIAMIntegration,
		),
		ReadContext: core.ReadResource(ReadAWSIAMIntegration),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "AWSIAMIntegrationUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/integrations/aws/iam/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
				NewResourceData: func() core.ResourceData {
					return &AWSIAMIntegrationWrapper{}
				},
			},
			ReadAWSIAMIntegration,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "AWSIAMIntegrationDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/integrations/aws/iam/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
			},
		),

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
