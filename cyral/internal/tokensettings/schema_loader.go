package tokensettings

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "tokensettings"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   "cyral_access_token_settings",
			Type:   core.DataSourceSchemaType,
			Schema: dataSourceSchema,
		},
		{
			Name:   "cyral_access_token_settings",
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
