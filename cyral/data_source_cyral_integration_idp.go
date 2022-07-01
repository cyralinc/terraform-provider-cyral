package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type IdPIntegrations struct {
	Connections *Connections `json:"connections,omitempty"`
}

type Connections struct {
	Connections []*Connection `json:"connections,omitempty"`
}

type Connection struct {
	DisplayName            string `json:"displayName,omitempty"`
	Alias                  string `json:"alias,omitempty"`
	SingleSignOnServiceURL string `json:"singleSignOnServiceURL,omitempty"`
	Enabled                bool   `json:"enabled,omitempty"`
}

func dataSourceIntegrationIdP() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter IdP integrations.",
		ReadContext: dataSourceIntegrationIdPRead,
		Schema: map[string]*schema.Schema{
			"display_name": {
				Description: "Filter results by the name of an existing IdP integration.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"type": {
				Description: "Filter results by the IdP integration type.",
				Optional:    true,
				Type:        schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"aad",
					"adfs",
					"forgerock",
					"gsuite",
					"okta",
					"pingone",
				}, false),
			},
			"idp_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Description: "Display name used in the Cyral control plane.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"alias": {
							Description: "Internal alias (ID) for this integration.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"single_sign_on_service_url": {
							Description: "Single sign on service URL.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"enabled": {
							Description: "True if the IdP integration is enabled.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
				Computed:    true,
				Description: "List of existing IdP integrations for the given filter criteria.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceIntegrationIdPRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	log.Printf("[DEBUG] Init dataSourceIntegrationIdPRead")
	c := m.(*client.Client)

	idpTypeFilter := d.Get("type").(string)

	var url = fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
	if idpTypeFilter != "" {
		url = fmt.Sprintf("https://%s/v1/integrations/saml?identityProvider=%s",
			c.ControlPlane, idpTypeFilter)
	}

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to execute request to read idp integrations", err.Error())
	}

	var idpIntegrations = IdPIntegrations{}
	if err := json.Unmarshal(body, &idpIntegrations); err != nil {
		return createError("Unable to unmarshal idp integrations", err.Error())
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", idpIntegrations)

	var idpList []interface{}
	displayNameFilter := d.Get("display_name").(string)

	log.Printf("[DEBUG] display_name: %s", displayNameFilter)
	log.Printf("[DEBUG] type: %s", idpTypeFilter)
	if idpIntegrations.Connections != nil {
		for _, connection := range idpIntegrations.Connections.Connections {
			log.Printf("[DEBUG] Connection: %#v", connection)
			if connection != nil {
				// Skip in case filters are non-empty but
				if displayNameFilter != "" && displayNameFilter != connection.DisplayName {
					continue
				}

				log.Printf("[DEBUG] Add connection to idp_set: %#v", connection)
				idpList = append(idpList, map[string]interface{}{
					"display_name":               connection.DisplayName,
					"alias":                      connection.Alias,
					"single_sign_on_service_url": connection.SingleSignOnServiceURL,
					"enabled":                    connection.Enabled,
				})
			}
		}
	}
	sort.Slice(idpList, func(p, q int) bool {
		return idpList[p].(map[string]interface{})["display_name"].(string) <
			idpList[q].(map[string]interface{})["display_name"].(string)
	})
	d.SetId(uuid.New().String())
	d.Set("idp_list", idpList)

	log.Printf("[DEBUG] End dataSourceIntegrationIdPRead")

	return nil
}
