package binding

import (
	"context"
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                  resourceName,
	ResourceType:                  resourcetype.Resource,
	SchemaReaderFactory:           func() core.SchemaReader { return &CreateBindingRequest{} },
	SchemaWriterFactoryGetMethod:  func(_ *schema.ResourceData) core.SchemaWriter { return &GetBindingResponse{} },
	SchemaWriterFactoryPostMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &CreateBindingResponse{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings",
			c.ControlPlane,
			d.Get(utils.SidecarIDKey).(string))
	},
	IdBasedURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
			c.ControlPlane,
			d.Get(utils.SidecarIDKey).(string),
			d.Get(utils.BindingIDKey).(string),
		)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [cyral repository to sidecar bindings](https://cyral.com/docs/sidecars/sidecar-assign-repo).",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		SchemaVersion: 2,
		Schema: map[string]*schema.Schema{
			utils.BindingIDKey: {
				Description: "ID of the binding. Computed and assigned to binding at the time of creation.",
				Computed:    true,
				Type:        schema.TypeString,
			},
			utils.SidecarIDKey: {
				Description: "ID of the sidecar that will be bound to the given repository.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			utils.RepositoryIDKey: {
				Description: "ID of the repository that will be bound to the sidecar.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			BindingEnabledKey: {
				Description: "Enable or disable all listener bindings.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			ListenerBindingKey: {
				Description: "The configuration for listeners associated with the binding. At least one `listener_binding` is required.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						utils.ListenerIDKey: {
							Description: "The sidecar listener that this binding is associated with.",
							Required:    true,
							Type:        schema.TypeString,
						},

						NodeIndexKey: {
							Description: "The index of the repo node that this binding is associated with.",
							Optional:    true,
							Type:        schema.TypeInt,
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
				ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				d.Set(utils.SidecarIDKey, ids[0])
				d.Set(utils.BindingIDKey, ids[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
