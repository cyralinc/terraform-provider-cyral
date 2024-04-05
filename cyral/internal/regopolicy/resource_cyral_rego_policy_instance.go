package regopolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Schema keys
	RegoPolicyInstanceResourceIDKey  = "id"
	RegoPolicyInstancePolicyIDKey    = "policy_id"
	RegoPolicyInstanceCategoryKey    = "category"
	RegoPolicyInstanceNameKey        = "name"
	RegoPolicyInstanceDescriptionKey = "description"
	RegoPolicyInstanceTemplateIDKey  = "template_id"
	RegoPolicyInstanceParametersKey  = "parameters"
	RegoPolicyInstanceEnabledKey     = "enabled"
	RegoPolicyInstanceScopeKey       = "scope"
	RegoPolicyInstanceRepoIDsKey     = "repo_ids"
	RegoPolicyInstanceTagsKey        = "tags"
	RegoPolicyInstanceDurationKey    = "duration"
	RegoPolicyInstanceLastUpdatedKey = "last_updated"
	RegoPolicyInstanceCreatedKey     = "created"
	RegoPolicyInstanceActorKey       = "actor"
	RegoPolicyInstanceActorTypeKey   = "actor_type"
	RegoPolicyInstanceTimestampKey   = "timestamp"
)

var (
	ReadRegoPolicyInstanceConfig = core.ResourceOperationConfig{
		ResourceName: "RegoPolicyInstanceRead",
		Type:         operationtype.Read,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/regopolicies/instances/%s/%s",
				c.ControlPlane,
				d.Get(RegoPolicyInstanceCategoryKey),
				d.Get(RegoPolicyInstancePolicyIDKey),
			)
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
			return &RegoPolicyInstance{}
		},
		RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Rego policy instance"},
	}

	regoPolicyChangeInformation = &schema.Resource{
		Schema: map[string]*schema.Schema{
			RegoPolicyInstanceActorKey: {
				Description: "Actor that performed the event.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			RegoPolicyInstanceActorTypeKey: {
				Description: "Actor type. Valid types are:" + utils.SupportedValuesAsMarkdown(actorTypes()),
				Type:        schema.TypeString,
				Computed:    true,
			},
			RegoPolicyInstanceTimestampKey: {
				Description: "Timestamp that the event happened.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
)

func ResourceRegoPolicyInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a [Rego Policy](https://cyral.com/docs/policy/rego-policy/overview#) instance." +
			"\n\n-> **Note** This resource can be used to create repo-level policies by specifying the repo IDs " +
			"associated to the policy `scope`. For more information, see the [scope](#nestedblock--scope) field.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName: "RegoPolicyInstanceCreate",
				Type:         operationtype.Create,
				HttpMethod:   http.MethodPost,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s",
						c.ControlPlane,
						d.Get(RegoPolicyInstanceCategoryKey),
					)
				},
				SchemaReaderFactory: func() core.SchemaReader { return &RegoPolicyInstancePayload{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &RegoPolicyInstanceKey{} },
			},
			ReadRegoPolicyInstanceConfig,
		),
		ReadContext: core.ReadResource(ReadRegoPolicyInstanceConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName: "RegoPolicyInstanceUpdate",
				Type:         operationtype.Update,
				HttpMethod:   http.MethodPut,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s/%s",
						c.ControlPlane,
						d.Get(RegoPolicyInstanceCategoryKey),
						d.Get(RegoPolicyInstancePolicyIDKey),
					)
				},
				SchemaReaderFactory: func() core.SchemaReader { return &RegoPolicyInstancePayload{} },
			},
			ReadRegoPolicyInstanceConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: "RegoPolicyInstanceDelete",
				Type:         operationtype.Delete,
				HttpMethod:   http.MethodDelete,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s/%s",
						c.ControlPlane,
						d.Get(RegoPolicyInstanceCategoryKey),
						d.Get(RegoPolicyInstancePolicyIDKey),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			RegoPolicyInstanceResourceIDKey: {
				Description: "The resource identifier. It is a composed ID that follows the format `{category}/{policy_id}`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			RegoPolicyInstancePolicyIDKey: {
				Description: "ID of this rego policy instance in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			RegoPolicyInstanceCategoryKey: {
				Description: "Policy category. List of supported categories:" +
					utils.SupportedValuesAsMarkdown(regoPolicyCategories()),
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(regoPolicyCategories(), false),
			},
			RegoPolicyInstanceNameKey: {
				Description: "Policy name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			RegoPolicyInstanceDescriptionKey: {
				Description: "Policy description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			RegoPolicyInstanceTemplateIDKey: {
				Description: "Policy template identifier. Predefined templates are:" +
					utils.SupportedValuesAsMarkdown(regoPolicyTemplateIDs()),
				Type:     schema.TypeString,
				Required: true,
			},
			RegoPolicyInstanceParametersKey: {
				Description: "Policy parameters. The parameters vary based on the policy template schema.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			RegoPolicyInstanceEnabledKey: {
				Description: "Enable/disable the policy. Defaults to `false` (Disabled).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			RegoPolicyInstanceScopeKey: {
				Description: fmt.Sprintf("Determines the scope that the policy applies to. It can be used to create "+
					"a repo-level policy by specifying the corresponding `%s` that this policy should be applied.",
					RegoPolicyInstanceRepoIDsKey),
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RegoPolicyInstanceRepoIDsKey: {
							Description: "A list of repository identifiers that belongs to the policy scope. The policy " +
								"will be applied at repo-level for every repository ID included in this list. This is equivalent " +
								"of creating a repo-level policy in the UI for a given repository.",
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			RegoPolicyInstanceTagsKey: {
				Description: "Tags that can be used to categorize the policy.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			RegoPolicyInstanceDurationKey: {
				Description: fmt.Sprintf(
					utils.DurationFieldDescriptionFmt,
					"Policy duration. The policy expires after the duration specified",
				),
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidationDurationString,
			},
			RegoPolicyInstanceLastUpdatedKey: {
				Description: "Information regarding the policy last update.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        regoPolicyChangeInformation,
			},
			RegoPolicyInstanceCreatedKey: {
				Description: "Information regarding the policy creation.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        regoPolicyChangeInformation,
			},
		},

		CustomizeDiff: func(ctx context.Context, resourceDiff *schema.ResourceDiff, i interface{}) error {
			computedKeysToChange := []string{RegoPolicyInstanceLastUpdatedKey}
			utils.SetKeysAsNewComputedIfPlanHasChanges(resourceDiff, computedKeysToChange)
			return nil
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				_ = d.Set(RegoPolicyInstanceCategoryKey, ids[0])
				_ = d.Set(RegoPolicyInstancePolicyIDKey, ids[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func regoPolicyCategories() []string {
	return []string{
		"SECURITY",
		"GRANT",
		"USER_DEFINED",
	}
}

func regoPolicyTemplateIDs() []string {
	return []string{
		"data-firewall",
		"data-masking",
		"data-protection",
		"EphemeralGrantPolicy",
		"rate-limit",
		"read-limit",
		"repository-protection",
		"service-account-abuse",
		"user-segmentation",
	}
}

func actorTypes() []string {
	return []string{
		"USER",
		"API_CLIENT",
	}
}
