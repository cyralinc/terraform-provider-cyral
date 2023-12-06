package samlcertificate

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SAMLCertificateData struct {
	Certificate string `json:"certificate,omitempty"`
}

func (data SAMLCertificateData) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set("certificate", data.Certificate)
	return nil
}
