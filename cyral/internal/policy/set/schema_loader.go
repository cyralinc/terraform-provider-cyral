package policyset

import "github.com/cyralinc/terraform-provider-cyral/cyral/core"

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "policyset"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{

		{
			Name:   policySetDataSourceName,
			Type:   core.DataSourceSchemaType,
			Schema: policySetDataSourceSchema,
		},
		{
			Name:   policySetResourceName,
			Type:   core.ResourceSchemaType,
			Schema: policySetResourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
