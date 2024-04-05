package repository

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Labels []string
type RepoNodes []*RepoNode

type RepoInfo struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	Host            string           `json:"repoHost"`
	Port            uint32           `json:"repoPort"`
	ConnParams      *ConnParams      `json:"connParams"`
	Labels          Labels           `json:"labels"`
	RepoNodes       RepoNodes        `json:"repoNodes,omitempty"`
	MongoDBSettings *MongoDBSettings `json:"mongoDbSettings,omitempty"`
}

type ConnParams struct {
	ConnDraining *ConnDraining `json:"connDraining"`
}

type ConnDraining struct {
	Auto     bool   `json:"auto"`
	WaitTime uint32 `json:"waitTime"`
}

type MongoDBSettings struct {
	ReplicaSetName string `json:"replicaSetName,omitempty"`
	ServerType     string `json:"serverType,omitempty"`
	SRVRecordName  string `json:"srvRecordName,omitempty"`
	Flavor         string `json:"flavor,omitempty"`
}

type RepoNode struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Dynamic bool   `json:"dynamic"`
}

type GetRepoByIDResponse struct {
	Repo RepoInfo `json:"repo"`
}

func (res *GetRepoByIDResponse) WriteToSchema(d *schema.ResourceData) error {
	return res.Repo.WriteToSchema(d)
}

func (res *RepoInfo) WriteToSchema(d *schema.ResourceData) error {
	d.Set(RepoTypeKey, res.Type)
	d.Set(RepoNameKey, res.Name)
	d.Set(RepoLabelsKey, res.Labels.AsInterface())
	d.Set(RepoConnDrainingKey, res.ConnParams.AsInterface())
	d.Set(RepoNodesKey, res.RepoNodes.AsInterface())
	d.Set(RepoMongoDBSettingsKey, res.MongoDBSettings.AsInterface())
	return nil
}

func (r *RepoInfo) ReadFromSchema(d *schema.ResourceData) error {
	r.ID = d.Id()
	r.Name = d.Get(RepoNameKey).(string)
	r.Type = d.Get(RepoTypeKey).(string)
	r.Labels = labelsFromInterface(d.Get(RepoLabelsKey).([]interface{}))
	r.RepoNodes = repoNodesFromInterface(d.Get(RepoNodesKey).([]interface{}))
	r.ConnParams = connDrainingFromInterface(d.Get(RepoConnDrainingKey).(*schema.Set).List())
	var mongoDBSettings = d.Get(RepoMongoDBSettingsKey).(*schema.Set).List()
	if r.Type == MongoDB && len(mongoDBSettings) == 0 {
		return fmt.Errorf("'%s' block must be provided when '%s=%s'", RepoMongoDBSettingsKey, utils.TypeKey, MongoDB)
	} else if r.Type != MongoDB && len(mongoDBSettings) > 0 {
		return fmt.Errorf("'%s' block is only allowed when '%s=%s'", RepoMongoDBSettingsKey, utils.TypeKey, MongoDB)
	}
	m, err := mongoDBSettingsFromInterface(mongoDBSettings)
	r.MongoDBSettings = m
	return err
}

func (l *Labels) AsInterface() []interface{} {
	if l == nil {
		return nil
	}
	labels := make([]interface{}, len(*l))
	for i, label := range *l {
		labels[i] = label
	}
	return labels
}

func labelsFromInterface(i []interface{}) Labels {
	labels := make([]string, len(i))
	for index, v := range i {
		labels[index] = v.(string)
	}
	return labels
}

func (c *ConnParams) AsInterface() []interface{} {
	if c == nil || c.ConnDraining == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		RepoConnDrainingAutoKey:     c.ConnDraining.Auto,
		RepoConnDrainingWaitTimeKey: c.ConnDraining.WaitTime,
	}}
}

func connDrainingFromInterface(i []interface{}) *ConnParams {
	if len(i) == 0 {
		return nil
	}
	return &ConnParams{
		ConnDraining: &ConnDraining{
			Auto:     i[0].(map[string]interface{})[RepoConnDrainingAutoKey].(bool),
			WaitTime: uint32(i[0].(map[string]interface{})[RepoConnDrainingWaitTimeKey].(int)),
		},
	}
}

func (r *RepoNodes) AsInterface() []interface{} {
	if r == nil {
		return nil
	}
	repoNodes := make([]interface{}, len(*r))
	for i, node := range *r {
		repoNodes[i] = map[string]interface{}{
			RepoNameKey:        node.Name,
			RepoHostKey:        node.Host,
			RepoPortKey:        node.Port,
			RepoNodeDynamicKey: node.Dynamic,
		}
	}
	return repoNodes
}

func repoNodesFromInterface(i []interface{}) RepoNodes {
	if len(i) == 0 {
		return nil
	}
	repoNodes := make(RepoNodes, len(i))
	for index, nodeInterface := range i {
		nodeMap := nodeInterface.(map[string]interface{})
		node := &RepoNode{
			Name:    nodeMap[RepoNameKey].(string),
			Host:    nodeMap[RepoHostKey].(string),
			Port:    uint32(nodeMap[RepoPortKey].(int)),
			Dynamic: nodeMap[RepoNodeDynamicKey].(bool),
		}
		repoNodes[index] = node
	}
	return repoNodes
}

func (m *MongoDBSettings) AsInterface() []interface{} {
	if m == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		RepoMongoDBReplicaSetNameKey: m.ReplicaSetName,
		RepoMongoDBServerTypeKey:     m.ServerType,
		RepoMongoDBSRVRecordName:     m.SRVRecordName,
		RepoMongoDBFlavorKey:         m.Flavor,
	}}
}

func mongoDBSettingsFromInterface(i []interface{}) (*MongoDBSettings, error) {
	if len(i) == 0 {
		return nil, nil
	}
	var replicaSetName = i[0].(map[string]interface{})[RepoMongoDBReplicaSetNameKey].(string)
	var serverType = i[0].(map[string]interface{})[RepoMongoDBServerTypeKey].(string)
	var srvRecordName = i[0].(map[string]interface{})[RepoMongoDBSRVRecordName].(string)
	if serverType == ReplicaSet && replicaSetName == "" {
		return nil, fmt.Errorf("'%s' must be provided when '%s=\"%s\"'", RepoMongoDBReplicaSetNameKey,
			RepoMongoDBServerTypeKey, ReplicaSet)
	}
	if serverType != ReplicaSet && replicaSetName != "" {
		return nil, fmt.Errorf("'%s' cannot be provided when '%s=\"%s\"'", RepoMongoDBReplicaSetNameKey,
			RepoMongoDBServerTypeKey, serverType)
	}
	if serverType == Standalone && srvRecordName != "" {
		return nil, fmt.Errorf(
			"'%s' cannot be provided when '%s=\"%s\"'",
			RepoMongoDBSRVRecordName,
			RepoMongoDBServerTypeKey,
			Standalone,
		)
	}
	return &MongoDBSettings{
		ReplicaSetName: i[0].(map[string]interface{})[RepoMongoDBReplicaSetNameKey].(string),
		ServerType:     i[0].(map[string]interface{})[RepoMongoDBServerTypeKey].(string),
		SRVRecordName:  i[0].(map[string]interface{})[RepoMongoDBSRVRecordName].(string),
		Flavor:         i[0].(map[string]interface{})[RepoMongoDBFlavorKey].(string),
	}, nil
}
