package health

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SidecarHealth struct {
	Status string `json:"status"`
}

func (health *SidecarHealth) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(utils.StatusKey, health.Status)

	return nil
}
