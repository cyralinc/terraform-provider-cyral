package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
					"okta",
					"adfs",
					"forgerock",
				}, false),
			},
			"idp_set": {
				Type: schema.TypeSet,
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
				Description: "Set of existing IdP integrations for the given filter criteria.",
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

	url := fmt.Sprintf("https://%s/v1/integrations/saml", c.ControlPlane)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to execute request to read idp integrations", err.Error())
	}

	var idpIntegrations = IdPIntegrations{}
	if err := json.Unmarshal(body, &idpIntegrations); err != nil {
		return createError("Unable to unmarshal idp integrations", err.Error())
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", idpIntegrations)

	var idpSet []interface{}
	displayName := d.Get("display_name").(string)
	idpType := d.Get("type").(string)

	log.Printf("[DEBUG] display_name: %s", displayName)
	log.Printf("[DEBUG] type: %s", idpType)
	if idpIntegrations.Connections != nil {
		for _, connection := range idpIntegrations.Connections.Connections {
			log.Printf("[DEBUG] Connection: %#v", connection)
			if connection != nil {
				// Conditions to return data:
				// 1. displayName is not empty and correspond to the one in the current connection;
				// 2. idpType is not empty and correspond to the first characters of the current connection alias;
				// 3. displayName and idpType are empty.
				if (displayName != "" && displayName == connection.DisplayName) ||
					(idpType != "" && idpType == connection.Alias[0:len(idpType)]) ||
					(displayName == "" && idpType == "") {
					log.Printf("[DEBUG] Add connection to idp_set: %#v", connection)
					idpSet = append(idpSet, map[string]interface{}{
						"display_name":               connection.DisplayName,
						"alias":                      connection.Alias,
						"single_sign_on_service_url": connection.SingleSignOnServiceURL,
						"enabled":                    connection.Enabled,
					})
				}
			}
		}
	}
	d.SetId(uuid.New().String())
	d.Set("idp_set", idpSet)

	log.Printf("[DEBUG] End dataSourceIntegrationIdPRead")

	return nil
}
