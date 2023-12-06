package tokensettings

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the access token settings. See also the data source " +
			"[`cyral_access_token_settings`](../data-source/access_token_settings.md)." +
			"\n\n-> **Note** The deletion of this terraform resource will reset the access " +
			"token settings to their corresponding default values.",
		CreateContext: core.CreateResource(updateConfig(), readConfig(resourcetype.Resource)),
		ReadContext:   core.ReadResource(readConfig(resourcetype.Resource)),
		UpdateContext: core.UpdateResource(updateConfig(), readConfig(resourcetype.Resource)),
		DeleteContext: core.DeleteResource(deleteConfig()),
		Schema:        getAccessTokenSettingsSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readConfig(rt resourcetype.ResourceType) core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
		Type:         operationtype.Read,
		ResourceType: rt,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/accessTokens/settings", c.ControlPlane)
		},
		SchemaWriterFactory: func(d *schema.ResourceData) core.SchemaWriter {
			return &AccessTokenSettings{}
		},
	}
}

func updateConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
		Type:         operationtype.Update,
		ResourceType: resourcetype.Resource,
		HttpMethod:   http.MethodPut,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/accessTokens/settings", c.ControlPlane)
		},
		SchemaReaderFactory: func() core.SchemaReader {
			return &AccessTokenSettings{}
		},
	}
}

// Since the access token settings resource is a global setting that is never deleted,
// the UpdateAccessTokenSettings API will be called here, with an empty body, so that
// the access token settings are reseted to their corresponding default values.
func deleteConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: resourceName,
		Type:         operationtype.Delete,
		ResourceType: resourcetype.Resource,
		HttpMethod:   http.MethodPut,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/accessTokens/settings", c.ControlPlane)
		},
	}
}
