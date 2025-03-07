package draft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

// Note that although the IdP type is SAML as seen from the control plane UI, we
// use GenericSAML in many variables, to make it more specific.
//
// The SAML IdP type is different from all others, in that it is _generic_. It
// allows the user to configure any type of IdP integration, as long as they
// provide the IdP metadata.
//
// This resource only covers the _draft_ of the integration. The draft is
// created on the control plane, and the SP metadata is generated. To complete
// the draft, one must perform the following steps:
//
// 1. Provide the SP metadata generated by this resource to the IdP.
//
// 2. Provide the IdP metadata to the `cyral_integration_idp_saml` resource.
//

var resourceContextHandler = core.HTTPContextHandler{
	ResourceName:                  resourceName,
	ResourceType:                  resourcetype.Resource,
	SchemaReaderFactory:           func() core.SchemaReader { return &CreateGenericSAMLDraftRequest{} },
	SchemaWriterFactoryGetMethod:  func(_ *schema.ResourceData) core.SchemaWriter { return &GenericSAMLDraftResponse{} },
	SchemaWriterFactoryPostMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &GenericSAMLDraftResponse{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages SAML IdP integration drafts." +
			"\n\nSee also the remaining SAML-related resources and data sources.",
		CreateContext: resourceContextHandler.CreateContextCustomErrorHandling(&readGenericSAMLDraftErrorHandler{}, nil),
		ReadContext:   resourceContextHandler.ReadContextCustomErrorHandling(&readGenericSAMLDraftErrorHandler{}),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Schema: map[string]*schema.Schema{
			// All of the input arguments must force recreation of
			// the resource, because the API does not support
			// updates. If you try to use the Create API to do
			// updates, it will create a new SAML draft altogether,
			// generating a new ID etc.
			"display_name": {
				Description:  "Display name used in the Cyral control plane.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
			},
			"disable_idp_initiated_login": {
				Description: "Whether or not IdP-Initiated login should be disabled for this generic SAML integration draft. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"idp_type": {
				Description: "Identity provider type. The value provided can be used as a filter when retrieving SAML integrations. See data source `cyral_integration_idp_saml`.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "saml-provider",
			},
			"attributes": {
				Description: "SAML Attribute names for the identity attributes required by the Cyral SP. Each attribute name MUST be at least 3 characters long.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users first name (given name). Defaults to `firstName`.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      idpsaml.DefaultUserAttributeFirstName,
							ValidateFunc: utils.ValidationStringLenAtLeast(3),
						},
						"last_name": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users last name (family name). Defaults to `lastName`.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      idpsaml.DefaultUserAttributeLastName,
							ValidateFunc: utils.ValidationStringLenAtLeast(3),
						},
						"email": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users email address. Defaults to `email`.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      idpsaml.DefaultUserAttributeEmail,
							ValidateFunc: utils.ValidationStringLenAtLeast(3),
						},
						"groups": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users group membership in the IdP. Defaults to `memberOf`.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      idpsaml.DefaultUserAttributeGroups,
							ValidateFunc: utils.ValidationStringLenAtLeast(3),
						},
					},
				},
			},
			"sp_metadata": {
				Description: "The SP Metadata document describing the Cyral service provider for this integration.",
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use `service_provider_metadata.xml_document` instead. This will be removed in the next major version of the provider.",
			},
			"service_provider_metadata": {
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
			"id": {
				Description: "ID of this resource in the Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

type readGenericSAMLDraftErrorHandler struct {
}

func (h *readGenericSAMLDraftErrorHandler) HandleError(
	ctx context.Context,
	d *schema.ResourceData,
	c *client.Client,
	err error,
) error {
	httpError, ok := err.(*client.HttpError)
	if !ok || httpError.StatusCode != http.StatusNotFound {
		return err
	}
	tflog.Debug(ctx, "SAML draft not found. Checking if completed draft exists.")

	query := utils.UrlQuery(map[string]string{
		"includeCompletedDrafts": "true",
		"displayName":            d.Get("display_name").(string),
		"idpType":                d.Get("idp_type").(string),
	})
	url := fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts%s",
		c.ControlPlane, query)
	body, err := c.DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("unable to read completed drafts: %w", err)
	}

	resp := ListGenericSAMLDraftsResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("when reading completed drafts: "+
			"unable to unmarshall JSON: %w", err)
	}

	myID := d.Id()
	found := false
	for _, draft := range resp.Drafts {
		if draft.ID == myID {
			found = true
			break
		}
	}
	if !found {
		tflog.Debug(ctx, fmt.Sprintf("Completed draft with ID %q "+
			"not found. Triggering recreation.", myID))
		d.SetId("")
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Found completed draft with ID %q.", myID))
	}
	return nil
}
