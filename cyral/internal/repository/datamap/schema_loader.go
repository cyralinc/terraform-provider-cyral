package datamap

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
)

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "datamap"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   "cyral_repository_datamap",
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
