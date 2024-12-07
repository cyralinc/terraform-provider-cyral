package logging

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

var resourceContextHandler = core.HTTPContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &LoggingIntegration{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &LoggingIntegration{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/logging", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a logging integration that can be used to push logs from Cyral to the corresponding logging system (E.g.: AWS CloudWatch, Splunk, SumoLogic, etc).",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Schema:        getIntegrationLogsSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
