package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BindingEnabledKey   = "enabled"
	ListenerBindingsKey = "listener_binding"
	NodeIndexKey        = "node_index"
)

type BindingResource struct {
	BindingID        string             `json:"id,omitempty"`
	RepoId           string             `json:"repoId,omitempty"`
	Enabled          bool               `json:"enabled,omitempty"`
	ListenerBindings []*ListenerBinding `json:"listenerBindings,omitempty"`
}

type ListenerBinding struct {
	ListenerID string `json:"listenerId,omitempty"`
	NodeIndex  uint32 `json:"nodeIndex,omitempty"`
}

type CreateBindingResponse struct {
	BindingID string `json:"bindingId,omitempty"`
}

func (r *CreateBindingResponse) WriteToSchema(d *schema.ResourceData) error {
	d.Set(BindingIDKey, r.BindingID)
	d.SetId(marshalComposedID(
		[]string{
			d.Get(SidecarIDKey).(string),
			r.BindingID,
		}, "/"))
	return nil
}

func (r *BindingResource) WriteToSchema(d *schema.ResourceData) error {
	d.Set(BindingIDKey, r.BindingID)
	d.Set(BindingEnabledKey, r.Enabled)
	d.Set(RepositoryIDKey, r.RepoId)
	d.Set(ListenerBindingsKey, r.ListenerBindingsAsInterface())
	return nil
}

func (r *BindingResource) ReadFromSchema(d *schema.ResourceData) error {
	r.Enabled = d.Get(BindingEnabledKey).(bool)
	r.RepoId = d.Get(RepositoryIDKey).(string)
	r.ListenerBindingsFromInterface(d.Get(ListenerBindingsKey).([]interface{}))
	return nil
}

func (r *BindingResource) ListenerBindingsAsInterface() []interface{} {
	if r.ListenerBindings == nil {
		return nil
	}
	listenerBindings := make([]interface{}, len(r.ListenerBindings))
	for i, listenerBinding := range r.ListenerBindings {
		listenerBindings[i] = map[string]interface{}{
			ListenerIDKey: listenerBinding.ListenerID,
			NodeIndexKey:  listenerBinding.NodeIndex,
		}
	}
	return listenerBindings
}

func (r *BindingResource) ListenerBindingsFromInterface(i []interface{}) {
	if len(i) == 0 {
		return
	}
	listenerBindings := make([]*ListenerBinding, len(i))
	for index, listenerBinding := range i {
		listenerBindings[index] = &ListenerBinding{
			ListenerID: listenerBinding.(map[string]interface{})[ListenerIDKey].(string),
			NodeIndex:  uint32(listenerBinding.(map[string]interface{})[NodeIndexKey].(int)),
		}
	}
	r.ListenerBindings = listenerBindings
}

var ReadRepositoryBindingConfig = ResourceOperationConfig{
	Name:       "RepositoryBindingResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
			c.ControlPlane,
			d.Get(SidecarIDKey).(string),
			d.Get(BindingIDKey).(string),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData {
		return &BindingResource{}
	},
}

func resourceRepositoryBinding() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [cyral repository to sidecar bindings](https://cyral.com/docs/sidecars/sidecar-assign-repo).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryBindingResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings",
						c.ControlPlane,
						d.Get(SidecarIDKey).(string))

				},
				NewResourceData: func() ResourceData {
					return &BindingResource{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &CreateBindingResponse{}
				},
			}, ReadSidecarListenersConfig,
		),
		ReadContext: ReadResource(ReadRepositoryBindingConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryBindingResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
						c.ControlPlane,
						d.Get(SidecarIDKey).(string),
						d.Get(BindingIDKey).(string),
					)

				},
				NewResourceData: func() ResourceData {
					return &BindingResource{}
				},
			}, ReadSidecarListenersConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryBindingResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
						c.ControlPlane,
						d.Get(SidecarIDKey).(string),
						d.Get(BindingIDKey).(string),
					)
				},
			},
		),

		SchemaVersion: 2,
		Schema: map[string]*schema.Schema{
			IDKey: {
				Description: "Terraform ID of this resource. Follows syntax `{sidecar_id}/{binding_id}`",
				Computed:    true,
				Type:        schema.TypeString,
			},
			BindingIDKey: {
				Description: "ID of the binding. Computed and assigned to binding at the time of creation.",
				Computed:    true,
				Type:        schema.TypeString,
			},
			SidecarIDKey: {
				Description: "ID of the sidecar that will be bound to the given repository.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			RepositoryIDKey: {
				Description: "ID of the repository that will be bound to the sidecar.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			BindingEnabledKey: {
				Description: "Enable or disable all listener bindings associated with this binding.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			ListenerBindingsKey: {
				Description: "The configuration for listeners associated with the binding. At least one `listener_binding` is required.",
				Type:        schema.TypeList,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ListenerIDKey: {
							Description: "The sidecar listener that this binding is associated with.",
							Optional:    true,
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
				ids, err := unmarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				d.Set(SidecarIDKey, ids[0])
				d.Set(BindingIDKey, ids[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
