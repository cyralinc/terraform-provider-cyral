package accessrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf("https://%s/v1/repos/%s/userAccounts/%s/accessRules",
		c.ControlPlane,
		d.Get("repository_id").(string),
		d.Get("user_account_id").(string),
	)
}

var readRepositoryAccessRulesConfig = core.ResourceOperationConfig{
	ResourceName: resourceName,
	Type:         operationtype.Read,
	HttpMethod:   http.MethodGet,
	URLFactory:   urlFactory,
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
		return &AccessRulesResponse{}
	},
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Repository access rule"},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manage access rules",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Create,
				HttpMethod:          http.MethodPut,
				URLFactory:          urlFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &AccessRulesResource{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &AccessRulesResponse{} },
			},
			readRepositoryAccessRulesConfig,
		),
		ReadContext: core.ReadResource(readRepositoryAccessRulesConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Update,
				HttpMethod:          http.MethodPut,
				URLFactory:          urlFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &AccessRulesResource{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &AccessRulesResponse{} },
			},
			readRepositoryAccessRulesConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: resourceName,
				Type:         operationtype.Delete,
				HttpMethod:   http.MethodDelete,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					// TODO Discuss why this is really necessary. We should be able
					// to use the same factory for all operations.
					idPieces, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
					if err != nil {
						panic(fmt.Sprintf("Failed to unmarshal access rules ID: %v", err))
					}
					repoID := idPieces[0]
					userAccountID := idPieces[1]

					return fmt.Sprintf("https://%s/v1/repos/%s/userAccounts/%s/accessRules",
						c.ControlPlane,
						repoID,
						userAccountID,
					)
				},
				RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: resourceName},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"repository_id": {
				Description: "ID of the repository.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"user_account_id": {
				Description: "ID of the database account. This should be the attribute `user_account_id` " +
					"of the resource `cyral_repository_user_account`.",
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},

			"rule": {
				Description: "An ordered list of access rules.",
				Required:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"identity": {
							Description: "The identity of the person/group getting access.",
							Required:    true,
							Type:        schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Identity type. List of supported values: " +
											utils.SupportedValuesAsMarkdown([]string{
												"username",
												"email",
												"group",
											}),
										Required: true,
										Type:     schema.TypeString,
									},

									"name": {
										Description: "The name of the person/group getting access.",
										Required:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MinItems: 1,
							MaxItems: 1,
						},

						"valid_from": {
							Description: "The start time for the grant. Format is: " +
								"`yyyy-mm-ddThh:mm:ssZ`. Eg. `2022-01-24T18:30:00Z`.",
							Optional: true,
							Type:     schema.TypeString,
						},

						"valid_until": {
							Description: "The end time for the grant. Format is: " +
								"`yyyy-mm-ddThh:mm:ssZ`. Eg. `2022-01-24T18:30:00Z`.",
							Optional: true,
							Type:     schema.TypeString,
						},

						"config": {
							Description: "Extra (optional) configuration parameters.",
							Optional:    true,
							MaxItems:    1,
							Type:        schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_ids": {
										Description: "Extra authorization policies, such as PagerDuty or DUO." +
											" Use the attribute `id` from resources `cyral_integration_pager_duty`" +
											" and `cyral_integration_mfa_duo`.",
										Required: true,
										Type:     schema.TypeList,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal ID: %v",
						err,
					)
				}
				err = d.Set("repository_id", ids[0])
				if err != nil {
					return nil, fmt.Errorf(
						"failed to set 'repository_id': %v",
						err,
					)
				}
				err = d.Set("user_account_id", ids[1])
				if err != nil {
					return nil, fmt.Errorf(
						"failed to set 'user_account_id: %v",
						err,
					)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
