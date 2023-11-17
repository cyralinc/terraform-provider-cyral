package tokensettings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "",
		CreateContext: core.CreateResource(updateConfig(), readConfig()),
		ReadContext:   core.ReadResource(readConfig()),
		UpdateContext: core.UpdateResource(updateConfig(), readConfig()),
		DeleteContext: resourceAccessTokenSettingsDelete,
		Schema:        getAccessTokenSettingsSchema(),
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

func resourceAccessTokenSettingsDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Since access token settings cannot be deleted, we just set the ID to
	// empty so that the resource can be removed from the terraform state.
	d.SetId("")
	return diag.Diagnostics{}
}
