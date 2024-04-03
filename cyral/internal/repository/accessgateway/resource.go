package accessgateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf(
		"https://%s/v1/repos/%s/accessGateway",
		c.ControlPlane,
		d.Get(utils.RepositoryIDKey).(string),
	)
}

var readRepositoryAccessGatewayConfig = core.ResourceOperationConfig{
	ResourceName: resourceName,
	Type:         operationtype.Read,
	HttpMethod:   http.MethodGet,
	URLFactory:   urlFactory,
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
		return &AccessGateway{}
	},
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: resourceName},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the sidecar and binding set as the access gateway for [cyral_repositories](./repositories.md).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Create,
				HttpMethod:          http.MethodPut,
				URLFactory:          urlFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &AccessGateway{} },
			},
			readRepositoryAccessGatewayConfig,
		),
		ReadContext: core.ReadResource(readRepositoryAccessGatewayConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Update,
				HttpMethod:          http.MethodPut,
				URLFactory:          urlFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &AccessGateway{} },
			},
			readRepositoryAccessGatewayConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Delete,
				HttpMethod:          http.MethodDelete,
				URLFactory:          urlFactory,
				RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: resourceName},
			},
		),

		Schema: map[string]*schema.Schema{
			utils.RepositoryIDKey: {
				Description: "ID of the repository the access gateway is associated with. This is also the " +
					"import ID for this resource.",
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			utils.SidecarIDKey: {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			utils.BindingIDKey: {
				Description: "ID of the binding that will be set as the access gateway for the given repository.  " +
					"Note that modifications to this field will result in terraform replacing the given " +
					"access gateway resource, since the access gateway must be deleted before binding. ",
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set(utils.RepositoryIDKey, d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
