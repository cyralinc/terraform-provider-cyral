package wizard

import "github.com/cyralinc/terraform-provider-cyral/cyral/core"

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "policyset"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   policyWizardsDataSourceName,
			Type:   core.DataSourceSchemaType,
			Schema: policyWizardsDataSourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
