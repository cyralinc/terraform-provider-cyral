package cyral

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultUserAttributeFirstName = "firstName"
	defaultUserAttributeLastName  = "lastName"
	defaultUserAttributeEmail     = "email"
	defaultUserAttributeGroups    = "memberOf"
)

type GenericSAMLDraft struct {
	ID                       string                  `json:"id"`
	DisplayName              string                  `json:"displayName"`
	IdpType                  string                  `json:"idpType"`
	DisableIdPInitiatedLogin bool                    `json:"disableIdpInitiatedLogin"`
	SPMetadata               *GenericSAMLSPMetadata  `json:"spMetadata,omitempty"`
	Attributes               *RequiredUserAttributes `json:"requiredUserAttributes,omitempty"`
	Completed                bool                    `json:"completed"`
}

type GenericSAMLSPMetadata struct {
	XMLDocument string `json:"xmlDocument"`
}

type RequiredUserAttributes struct {
	FirstName UserAttribute `json:"firstName,omitempty"`
	LastName  UserAttribute `json:"lastName,omitempty"`
	Email     UserAttribute `json:"email,omitempty"`
	Groups    UserAttribute `json:"groups,omitempty"`
}

func NewRequiredUserAttributes(firstName, lastName, email, groups string) *RequiredUserAttributes {
	return &RequiredUserAttributes{
		FirstName: UserAttribute{
			Name: firstName,
		},
		LastName: UserAttribute{
			Name: lastName,
		},
		Email: UserAttribute{
			Name: email,
		},
		Groups: UserAttribute{
			Name: groups,
		},
	}
}

func (uatt *RequiredUserAttributes) WriteToSchema(d *schema.ResourceData) error {
	var attributes []interface{}
	attributes = append(attributes, map[string]interface{}{
		"first_name": uatt.FirstName.Name,
		"last_name":  uatt.LastName.Name,
		"email":      uatt.Email.Name,
		"groups":     uatt.Groups.Name,
	})
	return d.Set("attributes", attributes)
}

func RequiredUserAttributesFromSchema(d *schema.ResourceData) (*RequiredUserAttributes, error) {
	attributesSet := d.Get("attributes").(*schema.Set).List()
	if len(attributesSet) > 1 {
		return nil, fmt.Errorf("Expected 'attributes' to be a set with at "+
			"most one element, got %d elements", len(attributesSet))
	} else if len(attributesSet) > 0 {
		attributesMap := attributesSet[0].(map[string]interface{})
		return NewRequiredUserAttributes(
			attributesMap["first_name"].(string),
			attributesMap["last_name"].(string),
			attributesMap["email"].(string),
			attributesMap["groups"].(string),
		), nil
	}
	return nil, nil
}

type UserAttribute struct {
	Name string `json:"name"`
}

type GenericSAMLIntegration struct {
	ID            string                    `json:"id"`
	DisplayName   string                    `json:"displayName"`
	IdpType       string                    `json:"idpType"`
	Disabled      bool                      `json:"disabled"`
	IdpDescriptor *GenericSAMLIdpDescriptor `json:"idpDescriptor,omitempty"`
	SPMetadata    *GenericSAMLSPMetadata    `json:"spMetadata,omitempty"`
	Attributes    *RequiredUserAttributes   `json:"attributes"`
}

func (integ *GenericSAMLIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(integ.ID)
	if integ.IdpDescriptor != nil {
		if err := d.Set("single_sign_on_service_url",
			integ.IdpDescriptor.SingleSignOnServiceURL,
		); err != nil {
			return err
		}
	}
	return nil
}

func (integ *GenericSAMLIntegration) ReadFromSchema(d *schema.ResourceData) error {
	integ.ID = d.Id()
	integ.IdpDescriptor = &GenericSAMLIdpDescriptor{
		SingleSignOnServiceURL: d.Get("single_sign_on_service_url").(string),
	}
	return nil
}

type GenericSAMLIdpDescriptor struct {
	SingleSignOnServiceURL     string `json:"singleSignOnServiceURL"`
	SigningCertificate         string `json:"signingCertificate"`
	DisableForceAuthentication bool   `json:"disableForceAuthentication"`
	SingleLogoutServiceURL     string `json:"singleLogoutServiceURL"`
}

type GenericSAMLIdpMetadata struct {
	URL string `json:"url,omitempty"`
	XML string `json:"xml,omitempty"`
}
