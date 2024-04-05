package awsiam

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func (wrapper AWSIAMIntegrationWrapper) WriteToSchema(d *schema.ResourceData) error {
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
