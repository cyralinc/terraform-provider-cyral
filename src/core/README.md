# Cyral Provider Core

The `core` package was created in order to put together all the code that is responsible
for managing the Provider itself and how to provide reusable functions and abstractions
that can be reused by the resources and data sources.

## How to use to create new resources and data sources

There are some main types that must be used to create a new resources and data sources:
`SchemaRegister`, `ResourceData`, `ResponseData` and `ResourceOperationConfig`. See the
examples below how to create your own implementation. See the source code for a more
in-depth documentation.

### model.go

```
// model.go

type NewResource struct {
	Name               string                 `json:"name,omitempty"`
	Description        string                 `json:"description,omitempty"`
}


func (r *NewResource) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("description", r.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	d.SetId(r.Name)
	return nil
}

func (r *NewResource) ReadFromSchema(d *schema.ResourceData) error {
	r.Name = d.Get("name").(string)
	r.Description = d.Get("description").(string)
	return nil
}
```

### datasource.go

```
// datasource.go

func init() {
	sr := &core.SchemaRegister{
        // Data source name
		Name:   "cyral_datalabel",
        // this is a function that will return the actual schema
		Schema: ResourceSchema,
        // Schema type
		Type:   core.ResourceSchema,
	}
	cyral.RegisterToProvider(sr)
}


func DataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Some description.",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
            Name:       "NewResourceRead",
            HttpMethod: http.MethodGet,
            CreateURL: func(d *schema.ResourceData, c *client.Client) string {
                return fmt.Sprintf("https://%s/v1/NewResource/%s", c.ControlPlane, d.Get("name").(string))
            },
            NewResponseData: func(d *schema.ResourceData) core.ResponseData {
                return &NewResource{}
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

func init() {
	sr := &core.SchemaRegister{
        // Data source name
		Name:   "cyral_datalabel",
        // this is a function that will return the actual schema
		Schema: ResourceSchema,
        // Schema type
		Type:   core.ResourceSchema,
	}
	cyral.RegisterToProvider(sr)
}

func ResourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Some description.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
                Name:       "NewResourceResourceRead",
                HttpMethod: http.MethodPost,
                CreateURL: func(d *schema.ResourceData, c *client.Client) string {
                    return fmt.Sprintf("https://%s/v1/NewResource", c.ControlPlane)
                },
                NewResponseData: func(d *schema.ResourceData) core.ResponseData {
                    return &NewResource{}
                },
            }, ReadNewResourceConfig,
		),
		ReadContext: core.ReadResource(ReadNewResourceConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "NewResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/NewResource/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &NewResource{} },
			}, ReadNewResourceConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "NewResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/NewResource/%s", c.ControlPlane, d.Id())
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

var ReadNewResourceConfig = core.ResourceOperationConfig{
	Name:       "NewResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/NewResource/%s", c.ControlPlane, d.Id())
	},
	NewResponseData:     func(_ *schema.ResourceData) core.ResponseData { return &NewResource{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "NewResource"},
}

```
