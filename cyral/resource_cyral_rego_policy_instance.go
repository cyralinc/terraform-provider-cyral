package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Schema keys
	regoPolicyInstanceIDKey          = "id"
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
	ReadRegoPolicyInstanceConfig = ResourceOperationConfig{
		Name:       "RegoPolicyInstanceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/regopolicies/instances/%s/%s",
				c.ControlPlane,
				d.Get("category"),
				d.Id(),
			)
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData {
			return &RegoPolicyInstance{}
		},
		RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "Rego policy instance"},
	}

	regoPolicyChangeInformation = &schema.Resource{
		Schema: map[string]*schema.Schema{
			regoPolicyInstanceActorKey: {
				Description: "Actor that performed the event.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			regoPolicyInstanceActorTypeKey: {
				Description: "Actor type. Valid types are:" + supportedTypesMarkdown(actorTypes()),
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
			"\n\n-> **NOTE** This resource can be used to create repo-level policies by specifying the repo IDs " +
			"associated to the policy `scope`. For more information, see the [scope](#nestedblock--scope) field.",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RegoPolicyInstanceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s",
						c.ControlPlane,
						d.Get("category"),
					)
				},
				NewResourceData: func() ResourceData {
					return &RegoPolicyInstancePayload{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &RegoPolicyInstanceCreateResponse{}
				},
			},
			ReadRegoPolicyInstanceConfig,
		),
		ReadContext: ReadResource(ReadRegoPolicyInstanceConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RegoPolicyInstanceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s/%s",
						c.ControlPlane,
						d.Get("category"),
						d.Id(),
					)
				},
				NewResourceData: func() ResourceData {
					return &RegoPolicyInstancePayload{}
				},
			},
			ReadRepositoryConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RegoPolicyInstanceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/regopolicies/instances/%s/%s",
						c.ControlPlane,
						d.Get("category"),
						d.Id(),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			regoPolicyInstanceIDKey: {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			regoPolicyInstanceCategoryKey: {
				Description: "Policy category. List of supported categories:" +
					supportedTypesMarkdown(regoPolicyCategories()),
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
				Description: "Policy template identifier. Valid templates are:" +
					supportedTypesMarkdown(regoPolicyTemplateIDs()),
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(regoPolicyTemplateIDs(), false),
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
							Optional: true,
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
				Description: "Policy duration. The policy expires after the duration specified. Should follow the Go " +
					"duration string format, which corresponds to a sequence of decimal numbers and a unit suffix. For " +
					"example: `300ms`, `2h45m`, `100s`, etc",
				Type:     schema.TypeString,
				Optional: true,
			},
			regoPolicyInstanceLastUpdatedKey: {
				Description: "Information regarding the policy last update.",
				Type:        schema.TypeSet,
				Computed:    true,
				MaxItems:    1,
				Elem:        regoPolicyChangeInformation,
			},
			regoPolicyInstanceCreatedKey: {
				Description: "Information regarding the policy creation.",
				Type:        schema.TypeSet,
				Computed:    true,
				MaxItems:    1,
				Elem:        regoPolicyChangeInformation,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
