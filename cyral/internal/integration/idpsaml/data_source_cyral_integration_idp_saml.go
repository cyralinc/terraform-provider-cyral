package idpsaml

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

type ListGenericSAMLIdpsRequest struct {
	DisplayName string `json:"displayName"`
	IdpType     string `json:"idpType"`
}

type ListGenericSAMLIdpsResponse struct {
	IdentityProviders []GenericSAMLIntegration `json:"identityProviders"`
}

func (resp *ListGenericSAMLIdpsResponse) WriteToSchema(d *schema.ResourceData) error {
	var idpList []interface{}
	for _, idp := range resp.IdentityProviders {
		var idpDescriptor []interface{}
		var spMetadata []interface{}
		var attributes []interface{}
		if idp.IdpDescriptor != nil {
			idpDescriptor = append(idpDescriptor, map[string]interface{}{
				"single_sign_on_service_url":   idp.IdpDescriptor.SingleSignOnServiceURL,
				"signing_certificate":          idp.IdpDescriptor.SigningCertificate,
				"disable_force_authentication": idp.IdpDescriptor.DisableForceAuthentication,
				"single_logout_service_url":    idp.IdpDescriptor.SingleLogoutServiceURL,
			})
		}
		if idp.SPMetadata != nil {
			spMetadata = idp.SPMetadata.ToList()
		}
		if idp.Attributes != nil {
			attributes = append(attributes, map[string]interface{}{
				"first_name": idp.Attributes.FirstName.Name,
				"last_name":  idp.Attributes.LastName.Name,
				"email":      idp.Attributes.Email.Name,
				"groups":     idp.Attributes.Groups.Name,
			})
		}
		idpList = append(idpList, map[string]interface{}{
			"id":             idp.ID,
			"display_name":   idp.DisplayName,
			"idp_type":       idp.IdpType,
			"disabled":       idp.Disabled,
			"idp_descriptor": idpDescriptor,
			"sp_metadata":    spMetadata,
			"attributes":     attributes,
		})
	}
	if err := d.Set("idp_list", idpList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	return nil
}

func dataSourceIntegrationIdPSAMLReadConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "IntegrationIdPSAMLDataSourceRead",
		Type:         operationtype.Read,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			query := utils.UrlQuery(map[string]string{
				"displayName": d.Get("display_name").(string),
				"idpType":     d.Get("idp_type").(string),
			})
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso%s", c.ControlPlane, query)
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &ListGenericSAMLIdpsResponse{} },
	}
}

func DataSourceIntegrationIdPSAML() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter SAML IdP integrations.",
		ReadContext: core.ReadResource(dataSourceIntegrationIdPSAMLReadConfig()),
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
			"idp_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of existing SAML IdP integrations that match the given filter criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Display name used in the Cyral control plane.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"display_name": {
							Description: "Display name used in the Cyral control plane.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"idp_type": {
							Description: "Indicates which type of identity provider this SSO integration is associated with.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"disabled": {
							Description: "True if the IdP integration is disabled.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"idp_descriptor": {
							Description: "The configuration information required by the Cyral SP, provided by the IdP.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"single_sign_on_service_url": {
										Description: "The IdP’s Single Sign-on Service (SS0) URL, where Cyral SP will send SAML AuthnRequests via SAML-POST binding.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"signing_certificate": {
										Description: "The signing certificate used by the Cyral SP to validate signed SAML assertions sent by the IdP.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"disable_force_authentication": {
										Description: "Indicates whether the identity provider must authenticate the presenter directly rather than rely on a previous security context. Defaults to false.",
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"single_logout_service_url": {
										Description: "The IdP’s Single Log-out Service (SL0) URL, where Cyral will send SAML AuthnRequests via SAML-POST binding.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"sp_metadata": {
							Description: "The SP Metadata fields describing the Cyral service provider for this integration.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"xml_document": {
										Description: "SP SAML metadata XML document.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"url": {
										Description: "URL where the metadata document can be downloaded.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"entity_id": {
										Description: "Entity ID defined in th SAML Metadata XML.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"single_logout_url": {
										Description: "The single logout URL defined in the SAML Metadata XML (SLO).",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"assertion_consumer_services": {
										Description: "The Assertion Consumer Services defined in the SAML Metadata XML.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Description: "The Assertion Consumer Service URL.",
													Type:        schema.TypeString,
													Computed:    true,
												},
												"index": {
													Description: "The index for the Assertion Consumer Service.",
													Type:        schema.TypeInt,
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"attributes": {
							Description: "SAML Attribute names for the identity attributes required by the Cyral SP.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_name": {
										Description: "The name of the attribute in the incoming SAML assertion containing the users first name (given name).",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"last_name": {
										Description: "The name of the attribute in the incoming SAML assertion containing the users last name (family name).",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"email": {
										Description: "The name of the attribute in the incoming SAML assertion containing the users email address.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"groups": {
										Description: "The name of the attribute in the incoming SAML assertion containing the users group membership in the IdP.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
