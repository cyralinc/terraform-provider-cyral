package datalabel

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

type packageSchema[T any] struct {
}

func (p *packageSchema[T]) Name() string {
	return "datalabel"
}

func (p *packageSchema[T]) Schemas() []*core.SchemaDescriptor {
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

func PackageSchema() core.PackageSchema[DataLabel] {
	return &packageSchema[DataLabel]{}
}
