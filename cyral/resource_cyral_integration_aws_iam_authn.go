package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
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

var ReadAWSIAMIntegration = ResourceOperationConfig{
	Name:       "AWSIAMIntegrationRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/integrations/aws/iam/%s",
			c.ControlPlane,
			d.Id(),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData {
		return &AWSIAMIntegrationWrapper{}
	},
	RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "AWS IAM Integration"},
}

func resourceIntegrationAWSIAM() *schema.Resource {
	return &schema.Resource{
		Description: "Authenticate users based on AWS IAM credentials",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "AWSIAMIntegrationCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/aws/iam", c.ControlPlane)
				},
				NewResourceData: func() ResourceData {
					return &AWSIAMIntegrationWrapper{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &IDBasedResponse{}
				},
			},
			ReadAWSIAMIntegration,
		),
		ReadContext: ReadResource(ReadAWSIAMIntegration),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "AWSIAMIntegrationUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/integrations/aws/iam/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
				NewResourceData: func() ResourceData {
					return &AWSIAMIntegrationWrapper{}
				},
			},
			ReadAWSIAMIntegration,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
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
			"id": {
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
