package binding

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BindingEnabledKey  = "enabled"
	ListenerBindingKey = "listener_binding"
	NodeIndexKey       = "node_index"
)

type Binding struct {
	BindingID        string             `json:"id,omitempty"`
	RepoId           string             `json:"repoId,omitempty"`
	Enabled          bool               `json:"enabled,omitempty"`
	ListenerBindings []*ListenerBinding `json:"listenerBindings,omitempty"`
}

type ListenerBinding struct {
	ListenerID string `json:"listenerId,omitempty"`
	NodeIndex  uint32 `json:"nodeIndex,omitempty"`
}

type CreateBindingRequest struct {
	SidecarID string   `json:"sidecarId,omitempty"`
	Binding   *Binding `json:"binding,omitempty"`
}

type CreateBindingResponse struct {
	BindingID string `json:"bindingId,omitempty"`
}

type GetBindingResponse struct {
	Binding *Binding `json:"binding,omitempty"`
}

func (r *CreateBindingResponse) WriteToSchema(d *schema.ResourceData) error {
	d.Set(utils.BindingIDKey, r.BindingID)
	d.SetId(utils.MarshalComposedID(
		[]string{
			d.Get(utils.SidecarIDKey).(string),
			r.BindingID,
		}, "/"))
	return nil
}

func (r *GetBindingResponse) WriteToSchema(d *schema.ResourceData) error {
	return r.Binding.WriteToSchema(d)
}

func (r *Binding) WriteToSchema(d *schema.ResourceData) error {
	d.Set(utils.BindingIDKey, r.BindingID)
	d.Set(BindingEnabledKey, r.Enabled)
	d.Set(utils.RepositoryIDKey, r.RepoId)
	d.Set(ListenerBindingKey, r.ListenerBindingsAsInterface())
	return nil
}

func (r *CreateBindingRequest) ReadFromSchema(d *schema.ResourceData) error {
	r.SidecarID = d.Get(utils.SidecarIDKey).(string)
	r.Binding = &Binding{}
	return r.Binding.ReadFromSchema(d)
}

func (r *Binding) ReadFromSchema(d *schema.ResourceData) error {
	r.BindingID = d.Get(utils.BindingIDKey).(string)
	r.Enabled = d.Get(BindingEnabledKey).(bool)
	r.RepoId = d.Get(utils.RepositoryIDKey).(string)
	r.ListenerBindingsFromInterface(d.Get(ListenerBindingKey).([]interface{}))
	return nil
}

func (r *Binding) ListenerBindingsAsInterface() []interface{} {
	if r.ListenerBindings == nil {
		return nil
	}
	listenerBindings := make([]interface{}, len(r.ListenerBindings))
	for i, listenerBinding := range r.ListenerBindings {
		listenerBindings[i] = map[string]interface{}{
			utils.ListenerIDKey: listenerBinding.ListenerID,
			NodeIndexKey:        listenerBinding.NodeIndex,
		}
	}
	return listenerBindings
}

func (r *Binding) ListenerBindingsFromInterface(i []interface{}) {
	if len(i) == 0 {
		return
	}
	listenerBindings := make([]*ListenerBinding, len(i))
	for index, listenerBinding := range i {
		listenerBindings[index] = &ListenerBinding{
			ListenerID: listenerBinding.(map[string]interface{})[utils.ListenerIDKey].(string),
			NodeIndex:  uint32(listenerBinding.(map[string]interface{})[NodeIndexKey].(int)),
		}
	}
	r.ListenerBindings = listenerBindings
}

var ReadRepositoryBindingConfig = core.ResourceOperationConfig{
	ResourceName: "RepositoryBindingResourceRead",
	HttpMethod:   http.MethodGet,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
			c.ControlPlane,
			d.Get(utils.SidecarIDKey).(string),
			d.Get(utils.BindingIDKey).(string),
		)
	},
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
		return &GetBindingResponse{}
	},
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Repository binding"},
}

func ResourceRepositoryBinding() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [cyral repository to sidecar bindings](https://cyral.com/docs/sidecars/sidecar-assign-repo).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName: "RepositoryBindingResourceCreate",
				HttpMethod:   http.MethodPost,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings",
						c.ControlPlane,
						d.Get(utils.SidecarIDKey).(string))

				},
				SchemaReaderFactory: func() core.SchemaReader { return &CreateBindingRequest{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &CreateBindingResponse{} },
			}, ReadRepositoryBindingConfig,
		),
		ReadContext: core.ReadResource(ReadRepositoryBindingConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName: "RepositoryBindingResourceUpdate",
				HttpMethod:   http.MethodPut,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
						c.ControlPlane,
						d.Get(utils.SidecarIDKey).(string),
						d.Get(utils.BindingIDKey).(string),
					)

				},
				SchemaReaderFactory: func() core.SchemaReader { return &CreateBindingRequest{} },
			}, ReadRepositoryBindingConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: "RepositoryBindingResourceDelete",
				HttpMethod:   http.MethodDelete,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/bindings/%s",
						c.ControlPlane,
						d.Get(utils.SidecarIDKey).(string),
						d.Get(utils.BindingIDKey).(string),
					)
				},
			},
		),

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
