package idpsaml_draft

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

type CreateGenericSAMLDraftRequest struct {
	DisplayName              string                          `json:"displayName"`
	DisableIdPInitiatedLogin bool                            `json:"disableIdPInitiatedLogin"`
	IdpType                  string                          `json:"idpType,omitempty"`
	Attributes               *idpsaml.RequiredUserAttributes `json:"attributes,omitempty"`
}

func (req *CreateGenericSAMLDraftRequest) ReadFromSchema(d *schema.ResourceData) error {
	req.DisplayName = d.Get("display_name").(string)
	req.DisableIdPInitiatedLogin = d.Get("disable_idp_initiated_login").(bool)
	req.IdpType = d.Get("idp_type").(string)

	attributes, err := idpsaml.RequiredUserAttributesFromSchema(d)
	if err != nil {
		return err
	}
	req.Attributes = attributes

	return nil
}

type GenericSAMLDraftResponse struct {
	Draft idpsaml.GenericSAMLDraft `json:"draft"`
}

func (resp *GenericSAMLDraftResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(resp.Draft.ID)
	if err := d.Set("display_name", resp.Draft.DisplayName); err != nil {
		return err
	}
	if err := d.Set("disable_idp_initiated_login", resp.Draft.DisableIdPInitiatedLogin); err != nil {
		return err
	}
	if err := d.Set("sp_metadata", resp.Draft.SPMetadata.XMLDocument); err != nil {
		return err
	}
	if resp.Draft.SPMetadata != nil {
		if err := resp.Draft.SPMetadata.WriteToSchema(d); err != nil {
			return err
		}
	}
	if err := d.Set("idp_type", resp.Draft.IdpType); err != nil {
		return err
	}
	if resp.Draft.Attributes != nil && utils.TypeSetNonEmpty(d, "attributes") {
		if err := resp.Draft.Attributes.WriteToSchema(d); err != nil {
			return err
		}
	}
	return nil
}

type ListGenericSAMLDraftsResponse struct {
	Drafts []idpsaml.GenericSAMLDraft `json:"drafts"`
}
