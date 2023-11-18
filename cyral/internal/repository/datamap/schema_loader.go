package datamap

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

type packageSchema[T any] struct {
}

func (p *packageSchema[T]) Name() string {
	return "datamap"
}

func (p *packageSchema[T]) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   "cyral_repository_datamap",
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema[DataMap] {
	return &packageSchema[DataMap]{}
}
