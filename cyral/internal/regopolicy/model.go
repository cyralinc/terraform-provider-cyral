package regopolicy

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RegoPolicyInstancePayload struct {
	RegoPolicyInstance RegoPolicyInstance `json:"instance"`
	Duration           string             `json:"duration,omitempty"`
}

func (policy *RegoPolicyInstancePayload) ReadFromSchema(d *schema.ResourceData) error {
	policy.RegoPolicyInstance.Name = d.Get(RegoPolicyInstanceNameKey).(string)
	policy.RegoPolicyInstance.Description = d.Get(RegoPolicyInstanceDescriptionKey).(string)
	policy.RegoPolicyInstance.TemplateID = d.Get(RegoPolicyInstanceTemplateIDKey).(string)
	policy.RegoPolicyInstance.Parameters = d.Get(RegoPolicyInstanceParametersKey).(string)
	policy.RegoPolicyInstance.Enabled = d.Get(RegoPolicyInstanceEnabledKey).(bool)
	policy.RegoPolicyInstance.Scope = NewScopeFromInterface(d.Get(RegoPolicyInstanceScopeKey))
	policy.RegoPolicyInstance.TagsFromInterfaceList(d.Get(RegoPolicyInstanceTagsKey).([]any))
	policy.Duration = d.Get(RegoPolicyInstanceDurationKey).(string)
	return nil
}

type RegoPolicyInstance struct {
	Name        string                        `json:"name"`
	Description string                        `json:"description,omitempty"`
	TemplateID  string                        `json:"templateId"`
	Parameters  string                        `json:"parameters,omitempty"`
	Enabled     bool                          `json:"enabled,omitempty"`
	Scope       *RegoPolicyInstanceScope      `json:"scope,omitempty"`
	Tags        []string                      `json:"tags,omitempty"`
	LastUpdated *RegoPolicyInstanceChangeInfo `json:"lastUpdated,omitempty"`
	Created     *RegoPolicyInstanceChangeInfo `json:"created,omitempty"`
}

func (policy *RegoPolicyInstance) WriteToSchema(d *schema.ResourceData) error {
	d.Set(RegoPolicyInstanceNameKey, policy.Name)
	d.Set(RegoPolicyInstanceDescriptionKey, policy.Description)
	d.Set(RegoPolicyInstanceTemplateIDKey, policy.TemplateID)
	d.Set(RegoPolicyInstanceParametersKey, policy.Parameters)
	d.Set(RegoPolicyInstanceEnabledKey, policy.Enabled)
	d.Set(RegoPolicyInstanceScopeKey, policy.Scope.ToInterfaceList())
	d.Set(RegoPolicyInstanceTagsKey, policy.TagsToInterfaceList())
	d.Set(RegoPolicyInstanceLastUpdatedKey, policy.LastUpdated.ToInterfaceList())
	d.Set(RegoPolicyInstanceCreatedKey, policy.Created.ToInterfaceList())
	return nil
}

func (policy *RegoPolicyInstance) TagsToInterfaceList() []any {
	if policy.Tags == nil {
		return nil
	}
	tags := make([]any, len(policy.Tags))
	for index, tag := range policy.Tags {
		tags[index] = tag
	}
	return tags
}

func (policy *RegoPolicyInstance) TagsFromInterfaceList(tagsInterfaceList []any) {
	tags := make([]string, len(tagsInterfaceList))
	for index, tag := range tagsInterfaceList {
		tags[index] = tag.(string)
	}
	policy.Tags = tags
}

type RegoPolicyInstanceScope struct {
	RepoIDs []string `json:"repoIds,omitempty"`
}

func NewScopeFromInterface(scopeInterface any) *RegoPolicyInstanceScope {
	if scopeInterface == nil {
		return nil
	}
	scopeInterfaceList := scopeInterface.(*schema.Set).List()
	if len(scopeInterfaceList) == 0 {
		return nil
	}
	scopeMap := scopeInterfaceList[0].(map[string]any)
	repoIDsInterfaceList := scopeMap[RegoPolicyInstanceRepoIDsKey].([]any)
	repoIDs := make([]string, len(repoIDsInterfaceList))
	for index, repoID := range repoIDsInterfaceList {
		repoIDs[index] = repoID.(string)
	}
	return &RegoPolicyInstanceScope{
		RepoIDs: repoIDs,
	}
}

func (scope *RegoPolicyInstanceScope) ToInterfaceList() []any {
	if scope == nil {
		return nil
	}
	return []any{
		map[string]any{
			RegoPolicyInstanceRepoIDsKey: scope.RepoIDsToInterfaceList(),
		},
	}
}

func (scope *RegoPolicyInstanceScope) RepoIDsToInterfaceList() []any {
	if scope.RepoIDs == nil {
		return nil
	}
	repoIDs := make([]any, len(scope.RepoIDs))
	for index, repoID := range scope.RepoIDs {
		repoIDs[index] = repoID
	}
	return repoIDs
}

type RegoPolicyInstanceChangeInfo struct {
	Actor     string `json:"actor,omitempty"`
	ActorType string `json:"actorType,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

func (changeInfo *RegoPolicyInstanceChangeInfo) ToInterfaceList() []any {
	if changeInfo == nil {
		return nil
	}
	return []any{
		map[string]any{
			RegoPolicyInstanceActorKey:     changeInfo.Actor,
			RegoPolicyInstanceActorTypeKey: changeInfo.ActorType,
			RegoPolicyInstanceTimestampKey: changeInfo.Timestamp,
		},
	}
}

type RegoPolicyInstanceKey struct {
	Category string `json:"category"`
	ID       string `json:"id"`
}

func (key RegoPolicyInstanceKey) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(utils.MarshalComposedID([]string{key.Category, key.ID}, "/"))
	d.Set(RegoPolicyInstancePolicyIDKey, key.ID)
	return nil
}
