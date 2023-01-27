package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *InsertPolicyInstanceRequest) ReadFromSchema(d *schema.ResourceData) error {
	scope := &Scope{}
	for _, scopeObj := range d.Get("scope").([]interface{}) {
		scopeMap := scopeObj.(map[string]interface{})
		repoIds := scopeMap["repo_ids"]
		for _, repoId := range repoIds.([]interface{}) {
			scope.RepoIds = append(scope.RepoIds, repoId.(string))
		}
	}

	lastUpdated := &ChangeInfo{}
	for _, lastUpdatedObj := range d.Get("last_updated").([]interface{}) {
		lastUpdatedMap := lastUpdatedObj.(map[string]interface{})
		actor := lastUpdatedMap["actor"].(string)
		actorType := lastUpdatedMap["actor_type"].(int32)
		timestamp := lastUpdatedMap["timestamp"].(int64)
		lastUpdated.Actor = actor
		lastUpdated.ActorType = ChangeInfo_ActorType(actorType)
		lastUpdated.Timestamp = &timestamppb.Timestamp{Seconds: timestamp}
	}

	created := &ChangeInfo{}
	for _, createdObj := range d.Get("created").([]interface{}) {
		createdMap := createdObj.(map[string]interface{})
		actor := createdMap["actor"].(string)
		actorType := createdMap["actor_type"].(int32)
		timestamp := createdMap["timestamp"].(int64)
		created.Actor = actor
		created.ActorType = ChangeInfo_ActorType(actorType)
		created.Timestamp = &timestamppb.Timestamp{Seconds: timestamp}
	}

	duration := d.Get("duration").(int64)

	r.Category = d.Get("category").(Category)
	r.Data = PolicyInstanceDataRequest{
		Instance: &PolicyInstance{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			TemplateId:  d.Get("template_id").(string),
			Parameters:  d.Get("parameters").(string),
			Enabled:     d.Get("enabled").(bool),
			Scope:       scope,
			LastUpdated: lastUpdated,
			Created:     created,
		},
		Duration: &durationpb.Duration{Seconds: duration},
	}
	return nil
}

func (r *UpdatePolicyInstanceRequest) ReadFromSchema(d *schema.ResourceData) error {
	r.AccessGateway = &AccessGateway{
		BindingId: d.Get(BindingIDKey).(string),
		SidecarId: d.Get(SidecarIDKey).(string),
	}
	return nil
}

func (r *DeletePolicyInstanceRequest) ReadFromSchema(d *schema.ResourceData) error {
	r.AccessGateway = &AccessGateway{
		BindingId: d.Get(BindingIDKey).(string),
		SidecarId: d.Get(SidecarIDKey).(string),
	}
	return nil
}

func (r *ReadPolicyInstanceRequest) ReadFromSchema(d *schema.ResourceData) error {
	r.AccessGateway = &AccessGateway{
		BindingId: d.Get(BindingIDKey).(string),
		SidecarId: d.Get(SidecarIDKey).(string),
	}
	return nil
}

func (r *InsertPolicyInstanceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(RepositoryIDKey).(string))
	d.Set("name", r.instance.Name)
	d.Set("description", r.instance.Description)
	d.Set("template_id", r.instance.TemplateId)
	d.Set("parameters", r.instance.Parameters)
	d.Set("enabled", r.instance.Enabled)
	repoIds := r.instance.Scope.RepoIds
	scope := map[string]interface{}{
		"repo_ids": repoIds,
	}
	d.Set("scope", scope)
	d.Set("tags", r.instance.Tags)
	lastUpdated := map[string]interface{}{
		"actor":      r.instance.LastUpdated.Actor,
		"actor_type": r.instance.LastUpdated.ActorType,
		"timestamp":  r.instance.LastUpdated.Timestamp,
	}
	d.Set("last_updated", lastUpdated)
	created := map[string]interface{}{
		"actor":      r.instance.Created.Actor,
		"actor_type": r.instance.Created.ActorType,
		"timestamp":  r.instance.Created.Timestamp,
	}
	d.Set("created", created)
	return nil
}

func (r *UpdatePolicyInstanceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(RepositoryIDKey).(string))
	d.Set("name", r.instance.Name)
	d.Set("description", r.instance.Description)
	d.Set("template_id", r.instance.TemplateId)
	d.Set("parameters", r.instance.Parameters)
	d.Set("enabled", r.instance.Enabled)
	repoIds := r.instance.Scope.RepoIds
	scope := map[string]interface{}{
		"repo_ids": repoIds,
	}
	d.Set("scope", scope)
	d.Set("tags", r.instance.Tags)
	lastUpdated := map[string]interface{}{
		"actor":      r.instance.LastUpdated.Actor,
		"actor_type": r.instance.LastUpdated.ActorType,
		"timestamp":  r.instance.LastUpdated.Timestamp,
	}
	d.Set("last_updated", lastUpdated)
	created := map[string]interface{}{
		"actor":      r.instance.Created.Actor,
		"actor_type": r.instance.Created.ActorType,
		"timestamp":  r.instance.Created.Timestamp,
	}
	d.Set("created", created)
	return nil
}

func (r *DeletePolicyInstanceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(RepositoryIDKey).(string))
	d.Set("name", r.instance.Name)
	d.Set("description", r.instance.Description)
	d.Set("template_id", r.instance.TemplateId)
	d.Set("parameters", r.instance.Parameters)
	d.Set("enabled", r.instance.Enabled)
	repoIds := r.instance.Scope.RepoIds
	scope := map[string]interface{}{
		"repo_ids": repoIds,
	}
	d.Set("scope", scope)
	d.Set("tags", r.instance.Tags)
	lastUpdated := map[string]interface{}{
		"actor":      r.instance.LastUpdated.Actor,
		"actor_type": r.instance.LastUpdated.ActorType,
		"timestamp":  r.instance.LastUpdated.Timestamp,
	}
	d.Set("last_updated", lastUpdated)
	created := map[string]interface{}{
		"actor":      r.instance.Created.Actor,
		"actor_type": r.instance.Created.ActorType,
		"timestamp":  r.instance.Created.Timestamp,
	}
	d.Set("created", created)
	return nil
}

func (r *ReadPolicyInstanceResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(RepositoryIDKey).(string))
	d.Set("name", r.instance.Name)
	d.Set("description", r.instance.Description)
	d.Set("template_id", r.instance.TemplateId)
	d.Set("parameters", r.instance.Parameters)
	d.Set("enabled", r.instance.Enabled)
	repoIds := r.instance.Scope.RepoIds
	scope := map[string]interface{}{
		"repo_ids": repoIds,
	}
	d.Set("scope", scope)
	d.Set("tags", r.instance.Tags)
	lastUpdated := map[string]interface{}{
		"actor":      r.instance.LastUpdated.Actor,
		"actor_type": r.instance.LastUpdated.ActorType,
		"timestamp":  r.instance.LastUpdated.Timestamp,
	}
	d.Set("last_updated", lastUpdated)
	created := map[string]interface{}{
		"actor":      r.instance.Created.Actor,
		"actor_type": r.instance.Created.ActorType,
		"timestamp":  r.instance.Created.Timestamp,
	}
	d.Set("created", created)
	return nil
}

// Reading the instance
var ReadRegopolicyInstanceConfig = ResourceOperationConfig{
	Name:       "RegopolicyInstanceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		category := d.Get("category").(string)
		id := d.Get("regopolicy_id").(string)
		return fmt.Sprintf("https://%s/v1/regopolicies/instances/%s/%s", c.ControlPlane, category, id)
	},
	NewResponseData: func(d *schema.ResourceData) ResponseData {
		return &ReadPolicyInstanceResponse{}
	},
}

func resourceRegopolicyInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the regopolicy instances.",
		// Creating the instance
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RegopolicyInstanceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					category := d.Get("category").(string)
					return fmt.Sprintf("https://%s/v1/regopolicies/instances/%s/", c.ControlPlane, category)
				},
				// payload to pass
				NewResourceData: func() ResourceData {
					return &InsertPolicyInstanceRequest{}
				},
				// return from API
				NewResponseData: func(d *schema.ResourceData) ResponseData {
					nameFilter := d.Get("name").(string)
					if nameFilter == "" {
						return &GetDataLabelsResponse{}
					} else {
						return &GetDataLabelResponse{}
					}
				},
			},
			// reading to update terraform values
			ReadRegopolicyInstanceConfig,
		),
		// Read
		ReadContext: ReadResource(ReadRegopolicyInstanceConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RegopolicyInstanceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					category := d.Get("category").(string)
					return fmt.Sprintf("https://%s/v1/regopolicies/instances/%s/", c.ControlPlane, category)
				},
				// payload to pass
				NewResourceData: func() ResourceData {

				},
				// return from API
				NewResponseData: func(d *schema.ResourceData) ResponseData {
					nameFilter := d.Get("name").(string)
					if nameFilter == "" {
						return &GetDataLabelsResponse{}
					} else {
						return &GetDataLabelResponse{}
					}
				},
			},
			// reading to update terraform values
			ReadRegopolicyInstanceConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RegopolicyInstanceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					category := d.Get("category").(string)
					id := d.Get("id").(string)
					return fmt.Sprintf("https://%s/v1/regopolicies/instances/%s/%s", c.ControlPlane, category, id)
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"regopolicy_id": {
				Description: "ID of the repository the access gateway is associated with. This is also the " +
					"import ID for this resource.",
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"category": {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"template_id": {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"parameters": {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Whether the policy is enabled or not",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"scope": {
				Description: "Object that defines the scope of the policy, i.e. where it is applicable",
				Type:        schema.TypeSet,
				Computed:    true,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_ids": {
							Description: "List of repo ids where the policy is applicable",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"tags": {
				Description: "Tags used to categorize policy instance.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_updated": {
				Description: "Classification rules are used by the [Automatic Data Map](https://cyral.com/docs/policy/automatic-datamap) feature to automatically map data locations to labels. Currently, only `PREDEFINED` labels have classification rules.",
				Type:        schema.TypeSet,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actor": {
							Description: "Actor identification (e.g. email)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"actor_type": {
							Description: "Type of actor, if it is user or api",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"timestamp": {
							Description: "Timestamp for action of updating",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"created": {
				Description: "Classification rules are used by the [Automatic Data Map](https://cyral.com/docs/policy/automatic-datamap) feature to automatically map data locations to labels. Currently, only `PREDEFINED` labels have classification rules.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actor": {
							Description: "Actor identification (e.g. email)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"actor_type": {
							Description: "Type of actor, if it is user or api",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"timestamp": {
							Description: "Timestamp for the action of creating",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"duration": {
				Description: "Duration of the policy instance.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set(RepositoryIDKey, d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
