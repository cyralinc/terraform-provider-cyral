package credentials

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSidecarCredentialsRequest struct {
	SidecarID string `json:"sidecarId"`
}

func (r *CreateSidecarCredentialsRequest) ReadFromSchema(d *schema.ResourceData) error {
	r.SidecarID = d.Get("sidecar_id").(string)
	return nil
}

type SidecarCredentialsData struct {
	SidecarID    string `json:"sidecarId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (r *SidecarCredentialsData) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("client_id", r.ClientID); err != nil {
		return fmt.Errorf("error setting 'client_id' field: %w", err)
	}
	if r.ClientSecret != "" {
		if err := d.Set("client_secret", r.ClientSecret); err != nil {
			return fmt.Errorf("error setting 'client_secret' field: %w", err)
		}
	}
	if err := d.Set("sidecar_id", r.SidecarID); err != nil {
		return fmt.Errorf("error setting 'sidecar_id' field: %w", err)
	}
	d.SetId(r.ClientID)
	return nil
}
