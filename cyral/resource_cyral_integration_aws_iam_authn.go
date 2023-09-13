package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type APIWrapper struct {
	Integration *AWSIAMIntegrationResource `json:"iamIntegration"`
}

type AWSIAMIntegrationResource struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IAMRoleARNs []string `json:"iamRoleARNs"`
}

func (wrapper *APIWrapper) WriteToSchema(d *schema.ResourceData) error {
	integration := wrapper.Integration

	d.SetId(integration.ID)

	if err := d.Set("name", integration.Name); err != nil {
		return fmt.Errorf("error setting 'name': %w", err)
	}

	if err := d.Set("description", integration.Description); err != nil {
		return fmt.Errorf("error setting 'description': %w", err)
	}

	if err := d.Set("arns", integration.IAMRoleARNs); err != nil {
		return fmt.Errorf("error setting 'arns': %w", err)
	}
	return nil
}

func (wrapper *APIWrapper) ReadFromSchema(d *schema.ResourceData) error {
	wrapper.Integration = &AWSIAMIntegrationResource{}

	wrapper.Integration.Name = d.Get("name").(string)
	wrapper.Integration.Description = d.Get("description").(string)

	arns := d.Get("arns").([]interface{})
	stringARNs := make([]string, 0, len(arns))
	for _, arn := range arns {
		stringARNs = append(stringARNs, arn.(string))
	}

	wrapper.Integration.IAMRoleARNs = stringARNs
	return nil
}

type CreateAWSIAMIntegrationResponse struct {
	ID string `json:"id"`
}

func (c *CreateAWSIAMIntegrationResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(c.ID)
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
		return &APIWrapper{}
	},
	RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "AWS IAM AuthN Integration"},
}

func resourceIntegrationAWSIAMAuthN() *schema.Resource {
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
					return &APIWrapper{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &CreateAWSIAMIntegrationResponse{}
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
					return &APIWrapper{}
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
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				i interface{},
			) ([]*schema.ResourceData, error) {
				id := d.Id()
				err := d.Set("id", id)
				if err != nil {
					return nil, fmt.Errorf("failed to set 'id': %w", err)
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Terraform ID of this resource",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"name": {
				Description: "The name of this AWS IAM Authentication integration",
				Required:    true,
				Type:        schema.TypeString,
			},

			"description": {
				Description: "Optional description of this integration",
				Optional:    true,
				Type:        schema.TypeString,
			},

			"arns": {
				Description: "List of role ARNs which will be used for authentication",
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
