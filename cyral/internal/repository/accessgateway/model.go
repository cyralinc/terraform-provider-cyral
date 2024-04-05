package accessgateway

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AGData struct {
	SidecarId string `json:"sidecarId,omitempty"`
	BindingId string `json:"bindingId,omitempty"`
}

type AccessGateway struct {
	AGData *AGData `json:"accessGateway,omitempty"`
}

func (r *AccessGateway) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(utils.RepositoryIDKey).(string))
	d.Set(utils.SidecarIDKey, r.AGData.SidecarId)
	d.Set(utils.BindingIDKey, r.AGData.BindingId)
	return nil
}

func (r *AccessGateway) ReadFromSchema(d *schema.ResourceData) error {
	r.AGData = &AGData{
		BindingId: d.Get(utils.BindingIDKey).(string),
		SidecarId: d.Get(utils.SidecarIDKey).(string),
	}
	return nil
}
