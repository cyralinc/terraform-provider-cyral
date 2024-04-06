# Cyral Provider Core

The `core` package was created in order to put together all the code that is responsible
for managing the Provider itself and to provide reusable functions and abstractions
for resources and data sources.

## How to use to create new resources and data sources

There are some main types that must be used to create new resources and data sources:
`SchemaDescriptor`, `PackageSchema`, `SchemaReader`, `SchemaWriter` and
`ResourceOperationConfig`. In a nutshell, these abstractions provide the means to
teach the provider how to:

- interact with the API;
- describe the feature as a Terraform resource/data source;
- perform the translation from API to Terraform schema and vice-versa.

Use the files below as examples to create your own implementation. It is advised that
you create a single package to group both the resource and data sources for a given
feature/category and that you follow the same naming convention for all the files
to simplify future code changes by adopting a single code convention.

### contants.go

```go
// contants.go
package newfeature

const (
	// The resource and data source names are identical in this example,
	// but this may not always hold true
	resourceName          = "cyral_new_feature"
	dataSourceName        = "cyral_new_feature"
)
```

### model.go

```go
// model.go
package newfeature

type NewFeature struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r NewFeature) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("description", r.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	if err := d.Set("name", r.Name); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	d.SetId(r.Name)
	return nil
}

func (r *NewFeature) ReadFromSchema(d *schema.ResourceData) error {
	r.Name = d.Get("name").(string)
	r.Description = d.Get("description").(string)
	return nil
}
```

### datasource.go

Use the `ReadUpdateDeleteURLFactory` to provide the URL factory to read the data source from the API.

```go
// datasource.go
package newfeature

var dsContextHandler = core.DefaultContextHandler{
	ResourceName:        dataSourceName,
	ResourceType:        resourcetype.DataSource,
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &NewFeature{} },
	ReadUpdateDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/NewFeature/%s", c.ControlPlane, d.Get("my_id_field").(string))
	},
}

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Some description.",
		ReadContext: dsContextHandler.ReadContext(),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Retrieve the unique label with this name, if it exists.",
				Type:        schema.TypeString,
				Optional:    true,
			},
            "description": {
                Description: "Description of the data source.",
                Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
```

### resource.go

```go
// resource.go
package newfeature

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:        resourceName,
	ResourceType:        resourcetype.Resource,
	SchemaReaderFactory: func() core.SchemaReader { return &NewFeature{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &NewFeature{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/NewFeature", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Some description.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "...",
				Type:        schema.TypeString,
				Optional:    true,
			},
            "description": {
                Description: "...",
                Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
```

### schema_loader.go

```go
// schema_loader.go
package newfeature

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "newfeature"
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
```

### provider/schema_loader.go

Edit the existing `cyral/provider/schema_loader.go` file and add your new package schema
to function `packagesSchemas` as follows:

```go
package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	...
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/newfeature"
)

func packagesSchemas() []core.PackageSchema {
	v := []core.PackageSchema{
		...,
		newfeature.PackageSchema(),
	}
	return v
}
```
