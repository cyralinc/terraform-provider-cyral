package tokensettings

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the access token settings. See also the data source " +
			"[`cyral_access_token_settings`](../data-source/access_token_settings.md)." +
			"\n\n-> **Note** The deletion of this terraform resource will reset the access " +
			"token settings to their corresponding default values.",
		CreateContext: core.CreateResource(updateConfig(), readConfig()),
		ReadContext:   core.ReadResource(readConfig()),
		UpdateContext: core.UpdateResource(updateConfig(), readConfig()),
		DeleteContext: core.DeleteResource(deleteConfig()),
		Schema:        getAccessTokenSettingsSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		Name:       "AccessTokenSettingsRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/accessTokens/settings", c.ControlPlane)
		},
		NewResponseData: func(d *schema.ResourceData) core.ResponseData {
			return &AccessTokenSettings{}
		},
	}
}

func updateConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		Name:       "AccessTokenSettingsUpdate",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/accessTokens/settings", c.ControlPlane)
		},
		NewResourceData: func() core.ResourceData {
			return &AccessTokenSettings{}
		},
	}
}

// Since the access token settings resource is a global setting that is never deleted,
// the UpdateAccessTokenSettings API will be called here, with an empty body, so that
// the access token settings are reseted to their corresponding default values.
func deleteConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		Name:       "AccessTokenSettingsDelete",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/accessTokens/settings", c.ControlPlane)
		},
	}
}
