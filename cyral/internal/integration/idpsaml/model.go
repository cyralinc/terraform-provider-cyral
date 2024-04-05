package idpsaml

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DefaultUserAttributeFirstName = "firstName"
	DefaultUserAttributeLastName  = "lastName"
	DefaultUserAttributeEmail     = "email"
	DefaultUserAttributeGroups    = "memberOf"
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

type AssertionConsumerService struct {
	Url   string `json:"url"`
	Index int32  `json:"index"`
}
type GenericSAMLSPMetadata struct {
	XMLDocument string `json:"xmlDocument"`
	Url         string `json:"url"`
	// Entity ID defined in th SAML Metadata XML
	EntityID string `json:"entityID"`
	// The single logout URL defined in the SAML Metadata XML (SL0)
	SingleLogoutURL string `json:"singleLogoutURL"`
	// An array with the Assertion Consumer Services defined in the SAML Metadata XML
	AssertionConsumerServices []*AssertionConsumerService `json:"assertionConsumerServices"`
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

func (spMetadataObj *GenericSAMLSPMetadata) WriteToSchema(d *schema.ResourceData) error {
	return d.Set("service_provider_metadata", spMetadataObj.ToList())
}

func (spMetadataObj *GenericSAMLSPMetadata) ToList() []any {
	var spMetadata []any
	if spMetadataObj != nil {
		assertionConsumerServices := []map[string]any{}
		if spMetadataObj.AssertionConsumerServices != nil {
			for _, assertionConsumerService := range spMetadataObj.AssertionConsumerServices {
				acs := make(map[string]any)
				acs["url"] = assertionConsumerService.Url
				acs["index"] = assertionConsumerService.Index
				assertionConsumerServices = append(assertionConsumerServices, acs)
			}
		}
		spMetadata = append(spMetadata, map[string]any{
			"xml_document":                spMetadataObj.XMLDocument,
			"url":                         spMetadataObj.Url,
			"entity_id":                   spMetadataObj.EntityID,
			"single_logout_url":           spMetadataObj.SingleLogoutURL,
			"assertion_consumer_services": assertionConsumerServices,
		})
	}

	return spMetadata
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
	SingleSignOnServiceURL     string `json:"singleSignOnServiceURL,omitempty"`
	SigningCertificate         string `json:"signingCertificate,omitempty"`
	DisableForceAuthentication bool   `json:"disableForceAuthentication,omitempty"`
	SingleLogoutServiceURL     string `json:"singleLogoutServiceURL,omitempty"`
}

type GenericSAMLIdpMetadata struct {
	URL string `json:"url,omitempty"`
	XML string `json:"xml,omitempty"`
}

type KeycloakProvider struct{}

type IdentityProviderData struct {
	Keycloak KeycloakProvider `json:"keycloakProvider"`
}

func (data IdentityProviderData) WriteToSchema(d *schema.ResourceData) error {
	return nil
}

func (data *IdentityProviderData) ReadFromSchema(d *schema.ResourceData) error {
	return nil
}
