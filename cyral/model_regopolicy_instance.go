package cyral

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Scope struct {
	RepoIds []string `json:"repoIds,omitempty"`
}

type ChangeInfo_ActorType int32

const (
	ChangeInfo_USER       ChangeInfo_ActorType = 0
	ChangeInfo_API_CLIENT ChangeInfo_ActorType = 1
)

type Category int32

const (
	Category_UNKNOWN      Category = 0
	Category_SECURITY     Category = 1
	Category_GRANT        Category = 2
	Category_USER_DEFINED Category = 3
)

type ChangeInfo struct {
	Actor     string                 `json:"actor,omitempty"`
	ActorType ChangeInfo_ActorType   `json:"actorType,omitempty"`
	Timestamp *timestamppb.Timestamp `json:"timestamp,omitempty"`
}

type Key struct {
	Id       string   `json:"id,omitempty"`
	Category Category `json:"category,omitempty"`
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
	Instance *PolicyInstance      `json:"policyInstance,omitempty"`
	Duration *durationpb.Duration `json:"duration,omitempty"`
}

type ListPolicyInstancePartial struct {
	Key        Key    `json:"key,omitempty"`
	Name       string `json:"name,omitempty"`
	TemplateId string `json:"templateId,omitempty"`
}

type InsertPolicyInstanceRequest struct {
	Category Category                  `json:"category,omitempty"`
	Data     PolicyInstanceDataRequest `json:"data,omitempty"`
}

type UpdatePolicyInstanceRequest struct {
	Key  Key                       `json:"key,omitempty"`
	Data PolicyInstanceDataRequest `json:"data,omitempty"`
}

type DeletePolicyInstanceRequest struct {
	Key Key `json:"key,omitempty"`
}

type ReadPolicyInstanceRequest struct {
	Key Key `json:"key,omitempty"`
}

type ReadPolicyInstancesRequest struct {
	Category Category `json:"category,omitempty"`
}

type InsertPolicyInstanceResponse struct {
	Key Key `json:"key,omitempty"`
}

type UpdatePolicyInstanceResponse struct {
}

type DeletePolicyInstanceResponse struct {
	instance PolicyInstance `json:"instance,omitempty"`
}

type ReadPolicyInstanceResponse struct {
	instance PolicyInstance `json:"instance,omitempty"`
}

type ReadPolicyInstancesResponse struct {
	instances []ListPolicyInstancePartial `json:"instances,omitempty"`
}
