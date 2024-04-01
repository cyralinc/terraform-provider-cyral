package confextension

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	authorizationPurpose = "authorization"
	builtinCategory      = "builtin"

	PagerDutyTemplateType = "pagerduty"
	DuoMFATemplateType    = "duoMfa"
)

type IntegrationConfExtension struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Parameters   string `json:"parameters,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
	Category     string `json:"category,omitempty"`
	TemplateType string `json:"templateType,omitempty"`
}

func NewIntegrationConfExtension(templateType string) *IntegrationConfExtension {
	return &IntegrationConfExtension{
		Purpose:      authorizationPurpose,
		Category:     builtinCategory,
		TemplateType: templateType,
	}
}

type IntegrationConfExtensionParameters struct {
	APIToken       string `json:"apiToken,omitempty"`
	IntegrationKey string `json:"integrationKey,omitempty"`
	SecretKey      string `json:"secretKey,omitempty"`
	APIHostname    string `json:"apiHostname,omitempty"`
}

func (data *IntegrationConfExtension) WriteToSchema(d *schema.ResourceData) error {
	d.Set("id", data.ID)
	d.Set("name", data.Name)

	var parameters IntegrationConfExtensionParameters
	err := json.Unmarshal([]byte(data.Parameters), &parameters)
	if err != nil {
		return fmt.Errorf("unable to unmarshal parameters: %w", err)
	}

	switch data.TemplateType {
	case PagerDutyTemplateType:
		d.Set("api_token", parameters.APIToken)
	case DuoMFATemplateType:
		d.Set("integration_key", parameters.IntegrationKey)
		d.Set("secret_key", parameters.SecretKey)
		d.Set("api_hostname", parameters.APIHostname)
	}

	return nil
}

func (data *IntegrationConfExtension) ReadFromSchema(d *schema.ResourceData) error {
	data.ID = d.Get("id").(string)
	data.Name = d.Get("name").(string)

	var parameters IntegrationConfExtensionParameters
	switch data.TemplateType {
	case PagerDutyTemplateType:
		parameters.APIToken = d.Get("api_token").(string)
	case DuoMFATemplateType:
		parameters.IntegrationKey = d.Get("integration_key").(string)
		parameters.SecretKey = d.Get("secret_key").(string)
		parameters.APIHostname = d.Get("api_hostname").(string)
	}

	parametersBytes, err := json.Marshal(parameters)
	if err != nil {
		return fmt.Errorf("unable to marshal parameters: %w", err)
	}
	data.Parameters = string(parametersBytes)

	return nil
}

func ConfExtensionIntegrationCreate(templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: fmt.Sprintf("%s_IntegrationResourceCreate", templateType),
		Type:         operationtype.Create,
		HttpMethod:   http.MethodPost,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/integrations/confExtensions/instances", c.ControlPlane,
			)
		},
		SchemaReaderFactory: func() core.SchemaReader {
			return NewIntegrationConfExtension(templateType)
		},
		SchemaWriterFactory: core.DefaultSchemaWriterFactory,
	}
}

func ConfExtensionIntegrationRead(templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: fmt.Sprintf("%s_IntegrationResourceRead", templateType),
		Type:         operationtype.Read,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/integrations/confExtensions/instances/authorization/%s",
				c.ControlPlane, d.Id(),
			)
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
			return NewIntegrationConfExtension(templateType)
		},
	}
}

func ConfExtensionIntegrationUpdate(templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: fmt.Sprintf("%s_IntegrationResourceUpdate", templateType),
		Type:         operationtype.Update,
		HttpMethod:   http.MethodPut,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/integrations/confExtensions/instances/%s", c.ControlPlane, d.Id(),
			)
		},
		SchemaReaderFactory: func() core.SchemaReader {
			return NewIntegrationConfExtension(templateType)
		},
	}
}

func ConfExtensionIntegrationDelete(templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: fmt.Sprintf("%s_IntegrationResourceDelete", templateType),
		Type:         operationtype.Delete,
		HttpMethod:   http.MethodDelete,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/integrations/confExtensions/instances/authorization/%s",
				c.ControlPlane, d.Id(),
			)
		},
	}
}
