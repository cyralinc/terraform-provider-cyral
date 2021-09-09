package cyral

import (
	"context"
	"fmt"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CertificateSchema string

func (data CertificateSchema) WriteToSchema(d *schema.ResourceData) {
	d.Set("certificate", string(data))
}

func dataSourceSAMLCertificate() *schema.Resource {
	return &schema.Resource{
		Description: "X.509 certificate for saml integration validation",
		ReadContext: func(_ context.Context, rd *schema.ResourceData, i interface{}) diag.Diagnostics {
			c := i.(*client.Client)
			resp, err := c.DoRequest(fmt.Sprintf("https://%s/deploy/saml/cert.pem", c.ControlPlane[:len(c.ControlPlane)-len(":8000")]), "GET", nil)
			if err != nil {
				return diag.FromErr(err)
			}

			rd.SetId("0")
			certificate := strings.Split(string(resp), "\n")[1]
			return diag.FromErr(rd.Set("certificate", certificate))
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the certificate used for signing saml requests",
			},
		},
	}
}
