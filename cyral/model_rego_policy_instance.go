package cyral

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RegoPolicyInstancePayload struct {
	RegoPolicyInstance RegoPolicyInstance `json:"instance"`
	Duration           string             `json:"duration,omitempty"`
}

func (policy *RegoPolicyInstancePayload) ReadFromSchema(d *schema.ResourceData) error {
	policy.RegoPolicyInstance.Name = d.Get(regoPolicyInstanceNameKey).(string)
	policy.RegoPolicyInstance.Description = d.Get(regoPolicyInstanceDescriptionKey).(string)
	policy.RegoPolicyInstance.TemplateID = d.Get(regoPolicyInstanceTemplateIDKey).(string)
	policy.RegoPolicyInstance.Parameters = d.Get(regoPolicyInstanceParametersKey).(string)
	policy.RegoPolicyInstance.Enabled = d.Get(regoPolicyInstanceEnabledKey).(bool)
	policy.RegoPolicyInstance.Scope = NewScopeFromInterface(d.Get(regoPolicyInstanceScopeKey))
	policy.RegoPolicyInstance.TagsFromInterfaceList(d.Get(regoPolicyInstanceTagsKey).([]any))
	policy.Duration = d.Get(regoPolicyInstanceDurationKey).(string)
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
	d.Set(regoPolicyInstanceNameKey, policy.Name)
	d.Set(regoPolicyInstanceDescriptionKey, policy.Description)
	d.Set(regoPolicyInstanceTemplateIDKey, policy.TemplateID)
	d.Set(regoPolicyInstanceParametersKey, policy.Parameters)
	d.Set(regoPolicyInstanceEnabledKey, policy.Enabled)
	d.Set(regoPolicyInstanceScopeKey, policy.Scope.ToInterfaceList())
	d.Set(regoPolicyInstanceTagsKey, policy.TagsToInterfaceList())
	d.Set(regoPolicyInstanceLastUpdatedKey, policy.LastUpdated.ToInterfaceList())
	d.Set(regoPolicyInstanceCreatedKey, policy.Created.ToInterfaceList())
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
	scopeMap := scopeInterfaceList[0].(map[string]any)
	repoIDsInterfaceList := scopeMap[regoPolicyInstanceRepoIDsKey].([]any)
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
			regoPolicyInstanceRepoIDsKey: scope.RepoIDsToInterfaceList(),
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
			regoPolicyInstanceActorKey:     changeInfo.Actor,
			regoPolicyInstanceActorTypeKey: changeInfo.ActorType,
			regoPolicyInstanceTimestampKey: changeInfo.Timestamp,
		},
	}
}
