package hcvault

import "github.com/cyralinc/terraform-provider-cyral/cyral/core"

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "HC Vault Integration"
}

func (p *packageSchema) Schemas() []*core.SchemaDescriptor {
	return []*core.SchemaDescriptor{
		{
			Name:   "cyral_integration_hc_vault",
			Type:   core.ResourceSchemaType,
			Schema: resourceSchema,
		},
	}
}

func PackageSchema() core.PackageSchema {
	return &packageSchema{}
}
