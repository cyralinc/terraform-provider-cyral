package cyral

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Scope struct {
	RepoIds []string `json:"repoIds,omitempty"`
}

const (
	actorUser      = "USER"
	actorApiClient = "API_CLIENT"
)

func actorTypes() []string {
	return []string{
		actorUser,
		actorApiClient,
	}
}

const (
	categoryTypeUnknown    = "UNKNOWN"
	categoryTypePredefined = "SECURITY"
	categoryTypeCustom     = "GRANT"
	categoryUserDefined    = "USER_DEFINED"
)

func categoryTypes() []string {
	return []string{
		categoryTypeUnknown,
		categoryTypePredefined,
		categoryTypeCustom,
		categoryUserDefined,
	}
}

type ChangeInfo struct {
	Actor     string                 `json:"actor,omitempty"`
	ActorType string                 `json:"actorType,omitempty"`
	Timestamp *timestamppb.Timestamp `json:"timestamp,omitempty"`
}

type Key struct {
	Id       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
}

func (pi *PolicyInstance) TagsAsInterface() []interface{} {
	var tagIfaces []interface{}
	for _, tag := range pi.Tags {
		tagIfaces = append(tagIfaces, tag)
	}
	return tagIfaces
}

type PolicyInstance struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	TemplateId  string      `json:"templateId,omitempty"`
	Parameters  string      `json:"parameters,omitempty"`
	Enabled     bool        `json:"enabled,omitempty"`
	Scope       *Scope      `json:"scope,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	LastUpdated *ChangeInfo `json:"lastUpdated,omitempty"`
	Created     *ChangeInfo `json:"created,omitempty"`
}

// used for 'data' in requests
type PolicyInstanceDataRequest struct {
	Instance *PolicyInstance `json:"instance,omitempty"`
	Duration string          `json:"duration,omitempty"`
}

type ListPolicyInstancePartial struct {
	Key        Key    `json:"key,omitempty"`
	Name       string `json:"name,omitempty"`
	TemplateId string `json:"templateId,omitempty"`
}

type InsertPolicyInstanceRequest PolicyInstanceDataRequest

type UpdatePolicyInstanceRequest PolicyInstanceDataRequest

type DeletePolicyInstanceRequest struct {
	Id       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
}

type ReadPolicyInstanceRequest struct {
	Id       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
}

type ReadPolicyInstancesRequest struct {
	Category string `json:"category,omitempty"`
}

type InsertPolicyInstanceResponse struct {
	Id       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
}

type UpdatePolicyInstanceResponse struct {
}

type DeletePolicyInstanceResponse struct {
	Instance PolicyInstance `json:"instance,omitempty"`
}

type ReadPolicyInstanceResponse PolicyInstance

type ReadPolicyInstancesResponse struct {
	Instances []ListPolicyInstancePartial `json:"instances,omitempty"`
}
