package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSAMLCertificate() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a X.509 certificate used for signing SAML requests.",
		ReadContext: ReadResource(ResourceOperationConfig{
			Name:       "dataSourceSAMLCertificateRead",
			HttpMethod: http.MethodGet,
			CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf("https://%s/v1/integrations/saml/rsa/cert", c.ControlPlane)
			},
			ResponseData: &SAMLCertificateData{},
		}),
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Computed ID for this resource (locally computed to be used in Terraform state).",
				Computed: true,
				Type: schema.TypeString,
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The X.509 certificate used for signing SAML requests.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

type SAMLCertificateData struct {
	Certificate string `json:"certificate,omitempty"`
}

func (data SAMLCertificateData) WriteToSchema(d *schema.ResourceData) {
	d.SetId(uuid.New().String())
	d.Set("certificate", data.Certificate)
}

func (data *SAMLCertificateData) ReadFromSchema(d *schema.ResourceData) {}
