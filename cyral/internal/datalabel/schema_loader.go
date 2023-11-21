package datalabel

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "datalabel"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   "cyral_datalabel",
			Type:   core.DataSourceSchemaType,
			Schema: dataSourceSchema,
		},
		{
			Name:   "cyral_datalabel",
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
