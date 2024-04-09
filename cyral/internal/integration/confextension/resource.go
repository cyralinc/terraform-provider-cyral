package confextension

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateResource(resourceName, templateType string) schema.CreateContextFunc {
	return core.CreateResource(
		create(resourceName, templateType),
		read(resourceName, templateType),
	)
}

func ReadResource(resourceName, templateType string) schema.ReadContextFunc {
	return core.ReadResource(read(resourceName, templateType))
}

func UpdateResource(resourceName, templateType string) schema.UpdateContextFunc {
	return core.UpdateResource(
		update(resourceName, templateType),
		read(resourceName, templateType),
	)
}

func DeleteResource(resourceName, templateType string) schema.DeleteContextFunc {
	return core.DeleteResource(delete(resourceName))
}

func create(resourceName, templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
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

func read(resourceName, templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
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
		RequestErrorHandler: &core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "not found for key",
			OperationType:  operationtype.Read,
		},
	}
}

func update(resourceName, templateType string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
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

func delete(resourceName string) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
		Type:         operationtype.Delete,
		HttpMethod:   http.MethodDelete,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/integrations/confExtensions/instances/authorization/%s",
				c.ControlPlane, d.Id(),
			)
		},
		RequestErrorHandler: &core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "not found for key",
			OperationType:  operationtype.Delete,
		},
	}
}
