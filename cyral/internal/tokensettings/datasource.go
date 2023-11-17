package tokensettings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "",
		ReadContext: core.ReadResource(readConfig()),
		Schema:      utils.ConvertSchemaFieldsToComputed(getAccessTokenSettingsSchema()),
	}
}
