package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/src/client"
	"github.com/cyralinc/terraform-provider-cyral/src/core"
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Schema keys
	regoPolicyInstanceResourceIDKey  = "id"
	regoPolicyInstancePolicyIDKey    = "policy_id"
	regoPolicyInstanceCategoryKey    = "category"
	regoPolicyInstanceNameKey        = "name"
	regoPolicyInstanceDescriptionKey = "description"
	regoPolicyInstanceTemplateIDKey  = "template_id"
	regoPolicyInstanceParametersKey  = "parameters"
	regoPolicyInstanceEnabledKey     = "enabled"
	regoPolicyInstanceScopeKey       = "scope"
	regoPolicyInstanceRepoIDsKey     = "repo_ids"
	regoPolicyInstanceTagsKey        = "tags"
	regoPolicyInstanceDurationKey    = "duration"
	regoPolicyInstanceLastUpdatedKey = "last_updated"
	regoPolicyInstanceCreatedKey     = "created"
	regoPolicyInstanceActorKey       = "actor"
	regoPolicyInstanceActorTypeKey   = "actor_type"
	regoPolicyInstanceTimestampKey   = "timestamp"
)

var (
	ReadRegoPolicyInstanceConfig = core.ResourceOperationConfig{
		Name:       "RegoPolicyInstanceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/regopolicies/instances/%s/%s",
				c.ControlPlane,
				d.Get(regoPolicyInstanceCategoryKey),
				d.Get(regoPolicyInstancePolicyIDKey),
			)
		},
		NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
			return &RegoPolicyInstance{}
		},
		RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Rego policy instance"},
	}

	regoPolicyChangeInformation = &schema.Resource{
		Schema: map[string]*schema.Schema{
			regoPolicyInstanceActorKey: {
				Description: "Actor that performed the event.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			regoPolicyInstanceActorTypeKey: {
				Description: "Actor type. Valid types are:" + utils.SupportedTypesMarkdown(actorTypes()),
				Type:        schema.TypeString,
				Computed:    true,
			},
			regoPolicyInstanceTimestampKey: {
				Description: "Timestamp that the event happened.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
)

func resourceRegoPolicyInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a [Rego Policy](https://cyral.com/docs/policy/rego-policy/overview#) instance." +
			"\n\n-> **Note** This resource can be used to create repo-level policies by specifying the repo IDs " +
			"associated to the policy `scope`. For more information, see the [scope](#nestedblock--scope) field.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "RegoPolicyInstanceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s",
						c.ControlPlane,
						d.Get(regoPolicyInstanceCategoryKey),
					)
				},
				NewResourceData: func() core.ResourceData {
					return &RegoPolicyInstancePayload{}
				},
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
					return &RegoPolicyInstanceKey{}
				},
			},
			ReadRegoPolicyInstanceConfig,
		),
		ReadContext: core.ReadResource(ReadRegoPolicyInstanceConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "RegoPolicyInstanceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s/%s",
						c.ControlPlane,
						d.Get(regoPolicyInstanceCategoryKey),
						d.Get(regoPolicyInstancePolicyIDKey),
					)
				},
				NewResourceData: func() core.ResourceData {
					return &RegoPolicyInstancePayload{}
				},
			},
			ReadRegoPolicyInstanceConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "RegoPolicyInstanceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s/%s",
						c.ControlPlane,
						d.Get(regoPolicyInstanceCategoryKey),
						d.Get(regoPolicyInstancePolicyIDKey),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			regoPolicyInstanceResourceIDKey: {
				Description: "The resource identifier. It is a composed ID that follows the format `{category}/{policy_id}`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			regoPolicyInstancePolicyIDKey: {
				Description: "ID of this rego policy instance in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			regoPolicyInstanceCategoryKey: {
				Description: "Policy category. List of supported categories:" +
					utils.SupportedTypesMarkdown(regoPolicyCategories()),
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(regoPolicyCategories(), false),
			},
			regoPolicyInstanceNameKey: {
				Description: "Policy name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			regoPolicyInstanceDescriptionKey: {
				Description: "Policy description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			regoPolicyInstanceTemplateIDKey: {
				Description: "Policy template identifier. Predefined templates are:" +
					utils.SupportedTypesMarkdown(regoPolicyTemplateIDs()),
				Type:     schema.TypeString,
				Required: true,
			},
			regoPolicyInstanceParametersKey: {
				Description: "Policy parameters. The parameters vary based on the policy template schema.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			regoPolicyInstanceEnabledKey: {
				Description: "Enable/disable the policy. Defaults to `false` (Disabled).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			regoPolicyInstanceScopeKey: {
				Description: fmt.Sprintf("Determines the scope that the policy applies to. It can be used to create "+
					"a repo-level policy by specifying the corresponding `%s` that this policy should be applied.",
					regoPolicyInstanceRepoIDsKey),
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						regoPolicyInstanceRepoIDsKey: {
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
			regoPolicyInstanceTagsKey: {
				Description: "Tags that can be used to categorize the policy.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			regoPolicyInstanceDurationKey: {
				Description: "Policy duration. The policy expires after the duration specified. Should follow the protobuf " +
					"duration string format, which corresponds to a sequence of decimal numbers suffixed by a 's' at the " +
					"end, representing the duration in seconds. For example: `300s`, `60s`, `10.50s` etc",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidationDurationString,
			},
			regoPolicyInstanceLastUpdatedKey: {
				Description: "Information regarding the policy last update.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        regoPolicyChangeInformation,
			},
			regoPolicyInstanceCreatedKey: {
				Description: "Information regarding the policy creation.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        regoPolicyChangeInformation,
			},
		},

		CustomizeDiff: func(ctx context.Context, resourceDiff *schema.ResourceDiff, i interface{}) error {
			computedKeysToChange := []string{regoPolicyInstanceLastUpdatedKey}
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
				_ = d.Set(regoPolicyInstanceCategoryKey, ids[0])
				_ = d.Set(regoPolicyInstancePolicyIDKey, ids[1])
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
