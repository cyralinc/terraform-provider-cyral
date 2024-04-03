package datamap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName: "DataMapResourceCreate",
				Type:         operationtype.Create,
				HttpMethod:   http.MethodPut,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/datamap",
						c.ControlPlane,
						d.Get("repository_id").(string))
				},
				SchemaReaderFactory: func() core.SchemaReader { return &DataMapRequest{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &DataMap{} },
			}, readDataMapConfig,
		),

		ReadContext: core.ReadResource(readDataMapConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName: "DataMapResourceUpdate",
				Type:         operationtype.Update,
				HttpMethod:   http.MethodPut,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/datamap",
						c.ControlPlane,
						d.Get("repository_id").(string))
				},
				SchemaReaderFactory: func() core.SchemaReader { return &DataMapRequest{} },
			}, readDataMapConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: "DataMapResourceDelete",
				Type:         operationtype.Delete,
				HttpMethod:   http.MethodDelete,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/datamap",
						c.ControlPlane,
						d.Get("repository_id").(string))
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Description: "ID of the repository for which to configure a data map.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"mapping": {
				Description: "Mapping of a label to a list of data locations (attributes).",
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Description: "Label given to the attributes in this mapping.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"attributes": {
							Description: "List containing the specific locations of the data within the repo, " +
								"following the pattern `{SCHEMA}.{TABLE}.{ATTRIBUTE}` (ex: " +
								"`[your_schema_name.your_table_name.your_attr_name]`).\n\n" +
								"-> When referencing data in Dremio repository, please include the complete " +
								"location in `attributes`, separating spaces by dots. For example, an attribute " +
								"`my_attr` from table `my_tbl` within space `inner_space` within space `outer_space` " +
								"would be referenced as `outer_space.inner_space.my_tbl.my_attr`. For more information, " +
								"please see the [Policy Guide](https://cyral.com/docs/reference/policy/).",
							Type:     schema.TypeList,
							Required: true,
							// TODO: this ForceNew propagates to the parent attribute `mapping`. Therefore, any
							// new mapping will force recreation. In the future, it would be good to use the
							// `v1/repos/{repoID}/datamap/labels/{label}/attributes/{attribute}` endpoint to
							// avoid unnecessary resource recreation. -aholmquist 2022-08-04
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set("repository_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

var readDataMapConfig = core.ResourceOperationConfig{
	ResourceName: "DataMapResourceRead",
	Type:         operationtype.Read,
	HttpMethod:   http.MethodGet,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/repos/%s/datamap",
			c.ControlPlane,
			d.Get("repository_id").(string))
	},
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &DataMap{} },
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Data Map"},
}
