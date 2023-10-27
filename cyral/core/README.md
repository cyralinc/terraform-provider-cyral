# Cyral Provider Core

The `core` package was created in order to put together all the code that is responsible
for managing the Provider itself and to provide reusable functions and abstractions
for resources and data sources.

## How to use to create new resources and data sources

There are some main types that must be used to create a new resources and data sources:
`SchemaDescriptor`, `PackageSchema`, `ResourceData`, `ResponseData` and
`ResourceOperationConfig`. See the examples below how to create your own implementation.
See the source code for a more in-depth documentation.

### model.go

```
// model.go
package new_feature

type NewFeature struct {
	Name               string                 `json:"name,omitempty"`
	Description        string                 `json:"description,omitempty"`
}

func (r *NewFeature) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("description", r.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
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

```
// datasource.go
package new_feature

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

```
// resource.go
package new_feature

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
				NewFeatureData: func() core.ResourceData { return &NewFeature{} },
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

```
// schema_loader.go
package new_feature

type packageSchema struct {
}

func (p *packageSchema) Name() string {
	return "new_feature"
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
