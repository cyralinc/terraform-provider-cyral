package cyral

import (
	"fmt"
	"net/http"

	//	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type ListSAMLIdpsResponse struct {
}

func (resp *ListSAMLIdpsResponse) WriteToSchema(d *schema.ResourceData) error {
	return nil
}

func dataSourceIntegrationIdPSAMLReadConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "IntegrationIdPSAMLDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("TODO")
		},
		NewResponseData: func() ResponseData { return &ListSAMLIdpsResponse{} },
	}
}

func dataSourceIntegrationIdPSAML() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter SAML IdP integrations.",
		ReadContext: ReadResource(dataSourceIntegrationIdPSAMLReadConfig()),
		Schema: map[string]*schema.Schema{
			"display_name": {
				Description: "Filter results by the display name (as seen in the control plane UI) of existing SAML IdP integrations.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"idp_type": {
				Description: "Filter results by the SAML IdP integration type.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
