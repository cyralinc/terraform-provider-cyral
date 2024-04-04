package confanalysis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO This resource is more complex than it should be due to the fact that a call to
// repo creation automatically creates the conf/auth and also the conf/analysis configurations.
// Our API should be refactored so these operations should happen separately.

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis",
		c.ControlPlane,
		d.Get("repository_id"),
	)
}

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &UserConfig{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &RepositoryConfAnalysisData{} },
	PostURLFactory:               urlFactory,
	GetPutDeleteURLFactory:       urlFactory,
}

var requestErrorHandler = &core.IgnoreNotFoundByMessage{MessageMatches: "Cannot find config data for repo"}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages Repository Analysis Configuration. This resource allows configuring " +
			"[Data Activity Logs](https://cyral.com/docs/data-repos/config/#data-activity-logs), " +
			"[Alerts](https://cyral.com/docs/data-repos/config/#alerts) and " +
			"[Policy Enforcement](https://cyral.com/docs/data-repos/config/#policy-enforcement) " +
			"settings for Data Repositories.",
		CreateContext: resourceRepositoryConfAnalysisCreate,
		ReadContext: resourceContextHandler.ReadContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "Cannot find config data for repo",
			OperationType:  operationtype.Read,
		}),
		UpdateContext: resourceContextHandler.UpdateContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "Cannot find config data for repo",
			OperationType:  operationtype.Update,
		}, nil),
		DeleteContext: resourceContextHandler.DeleteContextCustomErrorHandling(&core.IgnoreNotFoundByMessage{
			ResName:        resourceName,
			MessageMatches: "Cannot find config data for repo",
			OperationType:  operationtype.Delete,
		}),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: repositoryConfAnalysisResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: UpgradeRepositoryConfAnalysisV0,
			},
		},

		Schema: repositoryConfAnalysisResourceSchemaV0().Schema,

		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set("repository_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func repositoryConfAnalysisResourceSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"repository_id": {
				Description: "The ID of an existing data repository resource that will be configured.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"redact": {
				Description:  "Valid values are: `all`, `none` and `watched`. If set to `all` it will enable the redact of all literal values, `none` will disable it, and `watched` will only redact values from tracked fields set in the Datamap.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: client.ValidateRepositoryConfAnalysisRedact(),
			},
			"alert_on_violation": {
				Description: "If set to `true` it will enable alert on policy violations.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"disable_pre_configured_alerts": {
				Description: "If set to `true` it will *disable* preconfigured alerts.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"enable_data_masking": {
				Description: "If set to `true` it will allow policies to force the masking " +
					" of specified data fields in the results of queries. " +
					"[Learn more](https://cyral.com/docs/using-cyral/masking/).",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"block_on_violation": {
				Description: "If set to `true` it will enable query blocking in case of a " +
					"policy violation.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_filter_analysis": {
				Description: "If set to `true` it will *disable* filter analysis.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"enable_dataset_rewrites": {
				Description: "If set to `true` it will enable rewriting queries.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"comment_annotation_groups": {
				Description: "Valid values are: `identity`, `client`, `repo`, `sidecar`. The " +
					"default behavior is to set only the `identity` when this option is " +
					"enabled, but you can also opt to add the contents of `client`, `repo`, " +
					" `sidecar` logging blocks as query comments. " +
					" [Learn more](https://support.cyral.com/support/solutions/articles/44002218978).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: client.ValidateRepositoryConfAnalysisCommentAnnotationGroups(),
				},
			},
			"log_groups": {
				Description: "Responsible for configuring the Log Settings. Valid values are documented below. The `log_groups` list support the following values: " +
					"\n  - `everything` - Enables all the Log Settings." +
					"\n  - `dql` - Enables the `DQLs` setting for `all requests`." +
					"\n  - `dml` - Enables the `DMLs` setting for `all requests`." +
					"\n  - `ddl` - Enables the `DDLs` setting for `all requests`." +
					"\n  - `sensitive & dql` - Enables the `DQLs` setting for `logged fields`." +
					"\n  - `sensitive & dml` - Enables the `DMLs` setting for `logged fields`." +
					"\n  - `sensitive & ddl` - Enables the `DDLs` setting for `logged fields`." +
					"\n  - `privileged` - Enables the `Privileged commands` setting." +
					"\n  - `port-scan` - Enables the `Port scans` setting." +
					"\n  - `auth-failure` - Enables the `Authentication failures` setting." +
					"\n  - `full-table-scan` - Enables the `Full scans` setting." +
					"\n  - `violations` - Enables the `Policy violations` setting." +
					"\n  - `connections` - Enables the `Connection activity` setting." +
					"\n  - `sensitive` - Log all queries manipulating sensitive fields (watches)" +
					"\n  - `data-classification` - Log all queries whose response was automatically classified as sensitive (credit card numbers, emails and so on)." +
					"\n  - `audit` - Log `sensitive`, `DQLs`, `DDLs`, `DMLs` and `privileged`." +
					"\n  - `error` - Log analysis errors." +
					"\n  - `new-connections` - Log new connections." +
					"\n  - `closed-connections` - Log closed connections.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: client.ValidateRepositoryConfAnalysisLogSettings(),
				},
			},
		},
	}
}

// Previously, the ID for cyral_repository_conf_analysis had the format
// {repository_id}/ConfAnalysis. The goal of this state upgrade is to remove
// this suffix `ConfAnalysis`.
func UpgradeRepositoryConfAnalysisV0(
	_ context.Context,
	rawState map[string]interface{},
	_ interface{},
) (map[string]interface{}, error) {
	rawState["id"] = rawState["repository_id"]
	return rawState, nil
}

func resourceRepositoryConfAnalysisCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRepositoryConfAnalysisCreate")
	c := m.(*client.Client)
	httpMethod := http.MethodPost
	if confAnalysisAlreadyExists(ctx, c, d) {
		httpMethod = http.MethodPut
	}
	tflog.Debug(ctx, "End resourceRepositoryConfAnalysisCreate")
	return core.CreateResource(
		core.ResourceOperationConfig{
			ResourceName:        resourceName,
			Type:                operationtype.Create,
			HttpMethod:          httpMethod,
			URLFactory:          urlFactory,
			SchemaReaderFactory: func() core.SchemaReader { return &UserConfig{} },
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &RepositoryConfAnalysisData{} },
		},
		core.ResourceOperationConfig{
			ResourceName:        resourceName,
			ResourceType:        resourcetype.Resource,
			Type:                operationtype.Read,
			HttpMethod:          http.MethodGet,
			URLFactory:          urlFactory,
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &RepositoryConfAnalysisData{} },
			RequestErrorHandler: &core.IgnoreNotFoundByMessage{
				ResName:        resourceName,
				MessageMatches: "Cannot find config data for repo",
				OperationType:  operationtype.Read,
			},
		},
	)(ctx, d, m)
}

func confAnalysisAlreadyExists(ctx context.Context, c *client.Client, d *schema.ResourceData) bool {
	_, err := c.DoRequest(ctx, urlFactory(d, c), http.MethodGet, nil)
	// See TODO on the top of this file
	if err != nil {

		tflog.Debug(ctx, fmt.Sprintf("Unable to read Conf Analysis resource for repository %s: %v",
			d.Get("repository_id").(string), err))
		return false
	}
	return true
}
