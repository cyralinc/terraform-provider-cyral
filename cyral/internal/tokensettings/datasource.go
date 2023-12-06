package tokensettings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves the access token settings. See also the resource " +
			"[`cyral_access_token_settings`](../resources/access_token_settings.md).",
		ReadContext: core.ReadResource(readConfig(resourcetype.DataSource)),
		Schema:      getAccessTokenSettingsSchema(true),
	}
}
