package binding

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ListenerBindings []*ListenerBinding

type Binding struct {
	BindingID        string           `json:"id,omitempty"`
	RepoId           string           `json:"repoId,omitempty"`
	Enabled          bool             `json:"enabled,omitempty"`
	ListenerBindings ListenerBindings `json:"listenerBindings,omitempty"`
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
	d.Set(ListenerBindingKey, r.ListenerBindings.AsInterface())
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

func (r *ListenerBindings) AsInterface() []interface{} {
	if r == nil {
		return nil
	}
	listenerBindings := make([]interface{}, len(*r))
	for i, listenerBinding := range *r {
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
