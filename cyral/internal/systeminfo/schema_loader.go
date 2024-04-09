package systeminfo

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "systeminfo"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   dataSourceName,
			Type:   core.DataSourceSchemaType,
			Schema: dataSourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
