package policywizardv1

import "github.com/cyralinc/terraform-provider-cyral/cyral/core"

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "policywizardv1"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{

		{
			Name:   dataSourceName,
			Type:   core.DataSourceSchemaType,
			Schema: dataSourceSchema,
		},

		{
			Name:   resourceName,
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
