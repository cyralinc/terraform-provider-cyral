# Cyral Provider Core

The `core` package was created in order to put together all the code that is responsible
for managing the Provider itself and to provide reusable functions and abstractions
for resources and data sources.

## How to use to create new resources and data sources

There are some main types that must be used to create a new resources and data sources:
`SchemaDescriptor`, `PackageSchema`, `ResourceData`, `ResponseData` and
`ResourceOperationConfig`. In a nutshell, these abstractions provide the means to
teach the provider how to interact with the API, how to describe the feature as a
Terraform resource/data source and finally teach the provider how to perform the
translation from API to Terraform schema and vice-versa.

Use the files below as examples to create your own implementation. It is advised that
you follow the same naming convention for all the files to simplify future code changes.

### model.go

```go
// model.go
package newfeature

type NewFeature struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r *NewFeature) WriteToSchema(d *schema.ResourceData) error {
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

```go
// datasource.go
package newfeature

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Some description.",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
            Name:       "NewFeatureRead",
            HttpMethod: http.MethodGet,
            CreateURL: func(d *schema.ResourceData, c *client.Client) string {
                return fmt.Sprintf("https://%s/v1/NewFeature/%s", c.ControlPlane, d.Get("name").(string))
            },
            NewResponseData: func(d *schema.ResourceData) core.ResponseData {
                return &NewFeature{}
            },
        }),
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

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Some description.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
                Name:       "NewFeatureResourceRead",
                HttpMethod: http.MethodPost,
                CreateURL: func(d *schema.ResourceData, c *client.Client) string {
                    return fmt.Sprintf("https://%s/v1/NewFeature", c.ControlPlane)
                },
                NewResponseData: func(d *schema.ResourceData) core.ResponseData {
                    return &NewFeature{}
                },
            }, ReadNewFeatureConfig,
		),
		ReadContext: core.ReadResource(ReadNewFeatureConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "NewFeatureUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/NewFeature/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &NewFeature{} },
			}, ReadNewFeatureConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "NewFeatureDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/NewFeature/%s", c.ControlPlane, d.Id())
				},
			},
		),
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

var ReadNewFeatureConfig = core.ResourceOperationConfig{
	Name:       "NewFeatureRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/NewFeature/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) core.ResponseData { return &NewFeature{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "NewFeature"},
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
			Name:   "cyral_newfeature",
			Type:   core.DataSourceSchemaType,
			Schema: dataSourceSchema,
		},
		{
			Name:   "cyral_newfeature",
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
