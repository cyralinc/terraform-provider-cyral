package samlcertificate

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a X.509 certificate used for signing SAML requests." +
			"\n\nSee also the remaining SAML-related resources and data sources.",
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			ResourceName: dataSourceName,
			Type:         operationtype.Read,
			HttpMethod:   http.MethodGet,
			URLFactory: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf("https://%s/v1/integrations/saml/rsa/cert", c.ControlPlane)
			},
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &SAMLCertificateData{} },
		}),
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Computed ID for this data source (locally computed to be used in Terraform state).",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The X.509 certificate used for signing SAML requests.",
			},
		},
	}
}
